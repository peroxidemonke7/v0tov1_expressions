package convertexpressions

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Precompiled regexes used during trie operations.
var (
	metadataRegex      = regexp.MustCompile(`^([^\{]+)\{([^\}]+)\}$`)
	arrayNotationRegex = regexp.MustCompile(`^([^\[]+)(\[[^\]]+\])$`)
)

// TrieNode represents a node in the trie structure
type TrieNode struct {
	children      map[string]*TrieNode // Named children (exact matches)
	wildcardChild *TrieNode            // Single wildcard child for dynamic segments

	// Context-aware children: map of context key (e.g. step type) to sub-trie root
	// These are checked when a ConversionContext is provided during matching
	contextChildren map[string]*TrieNode

	// Node metadata
	alias      string // Alias for relative path entry (e.g., "stage", "step")
	id         string // Unique identifier for this node
	v1Name     string // v1 output name; "-" means skip in output, "" means use key
	isWildcard bool   // True for dynamic IDs (STAGE_ID, STEP_ID)
	isArray    bool   // True if this node represents an array

	// Terminal node data
	isEnd  bool   // True if this is a complete conversion rule endpoint
	target string // Replacement pattern if this is an end node
}

// Trie represents the trie structure for conversion rules
type Trie struct {
	root       *TrieNode
	nodeIndex  map[string]*TrieNode   // Quick lookup by node ID
	aliasIndex map[string][]*TrieNode // Quick lookup by alias (single node per alias)
}

// NewTrie creates a new trie instance
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[string]*TrieNode),
		},
		nodeIndex:  make(map[string]*TrieNode),
		aliasIndex: make(map[string][]*TrieNode),
	}
}

// AttachRulesAt attaches conversion rules to a node identified by ID
// Rules are relative to the node
func (t *Trie) AttachRulesAt(nodeID string, rules []ConversionRule) {
	node := t.nodeIndex[nodeID]
	if node == nil {
		fmt.Printf("Warning: node ID '%s' not found in trie\n", nodeID)
		return
	}

	for _, rule := range rules {
		t.insertFromNode(node, rule.From, rule.To)
	}
}

// AttachRulesWithContextAt attaches context-aware conversion rules to a node.
// These rules are only matched when the ConversionContext has a matching StepType.
// contextKey is typically the step type (e.g., "Run", "Http").
func (t *Trie) AttachRulesWithContextAt(nodeID string, contextKey string, rules []ConversionRule) {
	node := t.nodeIndex[nodeID]
	if node == nil {
		fmt.Printf("Warning: node ID '%s' not found in trie\n", nodeID)
		return
	}

	if node.contextChildren == nil {
		node.contextChildren = make(map[string]*TrieNode)
	}

	// Get or create context-specific sub-trie root
	contextRoot := node.contextChildren[contextKey]
	if contextRoot == nil {
		contextRoot = &TrieNode{
			children: make(map[string]*TrieNode),
		}
		node.contextChildren[contextKey] = contextRoot
	}

	for _, rule := range rules {
		t.insertFromNode(contextRoot, rule.From, rule.To)
	}
}

// parseNodeMetadata extracts node name and metadata from a pattern part
// Format: "nodeName{alias: env, id: node_id}" -> ("nodeName", map[string]string{"alias": "env", "id": "node_id"})
func parseNodeMetadata(part string) (string, map[string]string) {
	metadata := make(map[string]string)

	if matches := metadataRegex.FindStringSubmatch(part); matches != nil {
		nodeName := matches[1]
		metadataStr := matches[2]

		// Parse key-value pairs
		pairs := strings.Split(metadataStr, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				metadata[key] = value
			}
		}

		return nodeName, metadata
	}

	return part, metadata
}

func (t *Trie) insertFromNode(startNode *TrieNode, pattern, target string) {
	patternParts := strings.Split(pattern, ".")
	targetParts := strings.Split(target, ".")
	node := startNode

	for i, part := range patternParts {
		// Parse metadata from the part
		nodeName, metadata := parseNodeMetadata(part)

		// Get corresponding target part for v1Name
		var v1Name string
		if i < len(targetParts) {
			targetPart := targetParts[i]

			// If this is the last pattern part and there are more target parts,
			// join all remaining target parts as the v1Name
			if i == len(patternParts)-1 && len(targetParts) > len(patternParts) {
				remainingTarget := targetParts[i:]
				if arrayMatch := arrayNotationRegex.FindStringSubmatch(remainingTarget[0]); arrayMatch != nil {
					remainingTarget[0] = arrayMatch[1]
				}
				v1Name = strings.Join(remainingTarget, ".")
			} else {
				if arrayMatch := arrayNotationRegex.FindStringSubmatch(targetPart); arrayMatch != nil {
					v1Name = arrayMatch[1]
				} else {
					v1Name = targetPart
				}
			}
		} else {
			v1Name = "-"
		}

		// Handle array notation in pattern
		arrayMatch := arrayNotationRegex.FindStringSubmatch(nodeName)
		if arrayMatch != nil {
			baseName := arrayMatch[1]
			if node.children[baseName] == nil {
				node.children[baseName] = &TrieNode{
					children: make(map[string]*TrieNode),
					isArray:  true,
					v1Name:   v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.children[baseName], metadata)
			}
			node = node.children[baseName]
		} else if nodeName == "*" {
			if node.wildcardChild == nil {
				node.wildcardChild = &TrieNode{
					children:   make(map[string]*TrieNode),
					isWildcard: true,
					v1Name:     v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.wildcardChild, metadata)
			}
			node = node.wildcardChild
		} else {
			if node.children[nodeName] == nil {
				node.children[nodeName] = &TrieNode{
					children: make(map[string]*TrieNode),
					v1Name:   v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.children[nodeName], metadata)
			}
			node = node.children[nodeName]
		}
	}

	node.isEnd = true
}

// applyMetadata applies metadata map to a node
func (t *Trie) applyMetadata(node *TrieNode, metadata map[string]string) {
	if alias, ok := metadata["alias"]; ok {
		node.alias = alias
		t.aliasIndex[alias] = append(t.aliasIndex[alias], node)
	}
	if id, ok := metadata["id"]; ok {
		node.id = id
		t.nodeIndex[id] = node
	}
}

// pathPart represents a parsed path segment
type pathPart struct {
	name       string
	arrayIndex string // e.g., "[0]", "[1]"
}

// matchContext holds state during path matching
type matchContext struct {
	arrayIndices []string // Stack of array indices encountered
	v1Path       []string // Built v1 path segments

	// Step context tracking for lazy resolution
	lastStepID  string // Last step ID seen in path (captured from wildcard after "steps")
	inStepsPath bool   // True if we're currently inside a "steps" path segment

	// FQN mode tracking
	fqnAttempted bool // True if we've already attempted FQN conversion (prevents infinite recursion)
}

func (t *Trie) Match(path string, context *ConversionContext) (string, bool) {
	parts := t.parsePath(path)
	if len(parts) == 0 {
		return path, false
	}

	// Try alias-based matching on the first segment
	firstSegment := parts[0].name
	if aliasedNodes, exists := t.aliasIndex[firstSegment]; exists {
		var bestResult string
		var bestMatched bool
		bestScore := -1

		for _, aliasNode := range aliasedNodes {
			ctx := &matchContext{
				arrayIndices: []string{},
				v1Path:       []string{},
			}

			// FQN MODE: For "step." alias, replace v1Path with CurrentStepV1Path
			// This handles expressions like "step.spec.bucket" where "step" refers to the current step
			if firstSegment == "step" && context != nil && context.UseFQN && context.CurrentStepV1Path != "" {
				ctx.v1Path = strings.Split(context.CurrentStepV1Path, ".")
			} else if aliasNode.v1Name != "" && aliasNode.v1Name != "-" {
				if aliasNode.v1Name == "*" {
					ctx.v1Path = append(ctx.v1Path, parts[0].name)
				} else {
					ctx.v1Path = append(ctx.v1Path, aliasNode.v1Name)
				}
			}

			// Set inStepsPath if we're matching "steps" alias - this enables step ID capture
			if firstSegment == "steps" {
				ctx.inStepsPath = true
			}

			if result, matched := t.matchRecursive(aliasNode, parts, 1, ctx, context); matched {
				score := t.calculateMatchScore(aliasNode, parts, 1, context)
				if score > bestScore {
					bestScore = score
					bestResult = result
					bestMatched = true
				}
			}
		}

		if bestMatched {
			return bestResult, true
		}
	}

	// Fallback: try from root
	ctx := &matchContext{
		arrayIndices: []string{},
		v1Path:       []string{},
	}
	if result, matched := t.matchRecursive(t.root, parts, 0, ctx, context); matched {
		return result, true
	}

	return path, false
}

// calculateMatchScore returns the number of non-wildcard matches in a path
func (t *Trie) calculateMatchScore(node *TrieNode, parts []pathPart, index int, context *ConversionContext) int {
	score := 0
	currentNode := node

	for i := index; i < len(parts); i++ {
		part := parts[i]

		// Exact match (non-wildcard) gets a point
		if child, exists := currentNode.children[part.name]; exists {
			score++
			currentNode = child
		} else if currentNode.wildcardChild != nil {
			// Wildcard match, no points
			currentNode = currentNode.wildcardChild
		} else {
			// No match, stop counting
			break
		}
	}

	return score
}

// parsePath splits a path into parts, handling array indices and nested <+...> expressions.
func (t *Trie) parsePath(path string) []pathPart {
	segments := splitPathSegments(path)
	var parts []pathPart

	for _, seg := range segments {
		if arrayMatch := arrayNotationRegex.FindStringSubmatch(seg); arrayMatch != nil {
			parts = append(parts, pathPart{
				name:       arrayMatch[1],
				arrayIndex: arrayMatch[2],
			})
		} else {
			parts = append(parts, pathPart{name: seg})
		}
	}

	return parts
}

// hasResolvableContext returns true if the ConversionContext can resolve a step type.
func hasResolvableContext(convContext *ConversionContext, ctx *matchContext) bool {
	return convContext != nil && (convContext.StepType != "" || convContext.CurrentStepType != "" ||
		(ctx.lastStepID != "" && convContext.StepTypeMap != nil))
}

// tryMatchChild attempts to match remaining parts through a child node using both
// context-specific and general rules. The order depends on whether context is available.
// isSkipped indicates the child node suppresses its output (v1Name="-").
func (t *Trie) tryMatchChild(child *TrieNode, parts []pathPart, nextIndex int, ctx *matchContext, convContext *ConversionContext, isSkipped bool) (string, bool) {
	hasCtx := hasResolvableContext(convContext, ctx)

	if hasCtx {
		// Try context-specific rules first (lazy resolution happens in tryContextMatch)
		if result, matched := t.tryContextMatch(child, parts, nextIndex, ctx, convContext); matched {
			return result, true
		}
		// Fall back to general rules
		if result, matched := t.matchRecursive(child, parts, nextIndex, ctx, convContext); matched {
			if !isSkipped || !t.isSkippedNodePassthrough(result, ctx, parts, nextIndex) {
				return result, true
			}
		}
	} else {
		// No context: try general rules first, then deterministic context fallback
		if result, matched := t.matchRecursive(child, parts, nextIndex, ctx, convContext); matched {
			if !isSkipped || !t.isSkippedNodePassthrough(result, ctx, parts, nextIndex) {
				return result, true
			}
		}
		if result, matched := t.tryContextMatch(child, parts, nextIndex, ctx, convContext); matched {
			return result, true
		}
	}

	return "", false
}

func (t *Trie) matchRecursive(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	// Base case: consumed all parts
	if index == len(parts) {
		return strings.Join(ctx.v1Path, "."), true
	}

	currentPart := parts[index]

	if currentPart.arrayIndex != "" {
		ctx.arrayIndices = append(ctx.arrayIndices, currentPart.arrayIndex)
	}

	// Try exact match first
	if child, exists := node.children[currentPart.name]; exists {
		v1Segment := t.buildV1Segment(child, currentPart)
		if v1Segment != "" {
			ctx.v1Path = append(ctx.v1Path, v1Segment)
		}

		// Track "steps" path for step ID capture
		wasInStepsPath := ctx.inStepsPath
		if currentPart.name == "steps" {
			ctx.inStepsPath = true
		}

		isSkipped := child.v1Name == "-"

		if result, matched := t.tryMatchChild(child, parts, index+1, ctx, convContext, isSkipped); matched {
			return result, true
		}

		// Backtrack
		ctx.inStepsPath = wasInStepsPath
		if v1Segment != "" {
			ctx.v1Path = ctx.v1Path[:len(ctx.v1Path)-1]
		}
	}

	// Try wildcard match
	if node.wildcardChild != nil {
		v1Segment := currentPart.name
		if currentPart.arrayIndex != "" {
			v1Segment += currentPart.arrayIndex
		}

		ctx.v1Path = append(ctx.v1Path, v1Segment)

		// Capture step ID when matching wildcard after "steps"
		// This enables lazy step type resolution for step.spec expressions
		savedStepID := ctx.lastStepID
		savedV1Path := make([]string, len(ctx.v1Path))
		copy(savedV1Path, ctx.v1Path)

		if ctx.inStepsPath {
			ctx.lastStepID = currentPart.name
			ctx.inStepsPath = false // Reset after capturing

			// FQN MODE: When we've just captured a step ID and UseFQN is enabled,
			// replace ctx.v1Path with the step's v1 FQN base path.
			if convContext != nil && convContext.UseFQN {
				stepID := currentPart.name
				var v1BasePath string

				// Look up the step's v1 FQN base path from StepV1PathMap
				if convContext.StepV1PathMap != nil {
					v1BasePath = convContext.StepV1PathMap[stepID]
				}

				if v1BasePath != "" {
					// Replace v1Path with the FQN base path segments
					ctx.v1Path = strings.Split(v1BasePath, ".")
				}
			}
		}

		if result, matched := t.tryMatchChild(node.wildcardChild, parts, index+1, ctx, convContext, false); matched {
			return result, true
		}

		// Backtrack
		ctx.lastStepID = savedStepID
		ctx.v1Path = savedV1Path[:len(savedV1Path)-1]
	}

	// No match — return partial conversion if we have progress
	if len(ctx.v1Path) > 0 {
		remaining := t.buildRemainingPath(parts, index)
		return strings.Join(ctx.v1Path, ".") + "." + remaining, true
	}

	return "", false
}

// tryContextMatch attempts to match remaining parts through a node's context-specific sub-trie.
// Returns the converted path and true if a context match was found.
// When no context is provided, it tries all available context types and returns the first match
// deterministically (by sorting context keys alphabetically).
//
// LAZY STEP TYPE RESOLUTION:
// If convContext.StepType is empty but we have a step ID captured in matchContext,
// we resolve the step type from StepTypeMap before attempting context matching.
func (t *Trie) tryContextMatch(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	if node.contextChildren == nil {
		return "", false
	}

	// Lazy step type resolution: resolve from captured step ID or current step
	resolvedStepType := t.resolveStepType(ctx, convContext)

	// If we have a resolved step type, try that specific context
	if resolvedStepType != "" {
		contextRoot, exists := node.contextChildren[resolvedStepType]
		if !exists {
			return "", false
		}

		result, matched := t.tryContextSubtree(contextRoot, parts, index, ctx, convContext)
		if matched {
			return result, true
		}
		return "", false
	}

	// No context provided - try all available contexts and pick the one
	// with the most node matches against the v0 input path.
	// Ties are broken alphabetically by context key for determinism.
	var contextKeys []string
	for key := range node.contextChildren {
		contextKeys = append(contextKeys, key)
	}

	if len(contextKeys) == 0 {
		return "", false
	}

	// Sort to ensure deterministic tie-breaking
	sortStrings(contextKeys)

	var bestResult string
	bestScore := -1
	bestKey := ""

	for _, contextKey := range contextKeys {
		contextRoot := node.contextChildren[contextKey]
		result, matched := t.tryContextSubtree(contextRoot, parts, index, ctx, convContext)
		if matched {
			score := t.countContextNodeMatches(contextRoot, parts, index)
			if score > bestScore || (score == bestScore && (bestKey == "" || contextKey < bestKey)) {
				bestScore = score
				bestResult = result
				bestKey = contextKey
			}
		}
	}

	if bestScore >= 0 {
		return bestResult, true
	}

	return "", false
}

// isSkippedNodePassthrough checks whether a result from a skipped (v1Name="-") node's
// deeper recursion is just a passthrough — i.e., the skipped node didn't contribute
// any real transformation. This covers two cases:
//   - Base case: path ends at the skipped node, result == v1Path joined (nothing added)
//   - Partial fallback: remaining parts are passed through raw, result == v1Path + "." + raw remaining
func (t *Trie) isSkippedNodePassthrough(result string, ctx *matchContext, parts []pathPart, fromIndex int) bool {
	v1Joined := strings.Join(ctx.v1Path, ".")

	// Base case: path ended at the skipped node
	if fromIndex >= len(parts) {
		return result == v1Joined
	}

	// Partial fallback: remaining parts passed through unchanged
	remaining := t.buildRemainingPath(parts, fromIndex)
	if remaining == "" {
		return result == v1Joined
	}
	return result == v1Joined+"."+remaining
}

// countContextNodeMatches counts the number of exact (non-wildcard) node matches
// in a context subtree against the remaining input path parts.
// This is used to determine the best-match context when no context is provided.
func (t *Trie) countContextNodeMatches(contextRoot *TrieNode, parts []pathPart, startIndex int) int {
	score := 0
	currentNode := contextRoot

	for i := startIndex; i < len(parts); i++ {
		part := parts[i]

		if child, exists := currentNode.children[part.name]; exists {
			score++
			currentNode = child
		} else if currentNode.wildcardChild != nil {
			currentNode = currentNode.wildcardChild
		} else {
			break
		}
	}

	return score
}

// tryContextSubtree attempts to match through a specific context sub-trie
func (t *Trie) tryContextSubtree(contextRoot *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	// Try matching remaining parts through the context sub-trie
	contextCtx := &matchContext{
		arrayIndices: ctx.arrayIndices,
		v1Path:       make([]string, len(ctx.v1Path)),
	}
	copy(contextCtx.v1Path, ctx.v1Path)

	result, matched := t.matchRecursive(contextRoot, parts, index, contextCtx, convContext)
	if matched {
		// Verify this is a complete match, not a partial fallback
		// The partial match fallback in matchRecursive can return true even when
		// it didn't properly match through the context sub-trie, so we need to validate.
		// A partial fallback produces: v1Path + "." + remainingRaw (raw parts unchanged)
		// A real context match produces: v1Path + transformed parts from context sub-trie
		remainingRaw := t.buildRemainingPath(parts, index)
		if remainingRaw != "" {
			expectedPartialFallback := strings.Join(contextCtx.v1Path, ".") + "." + remainingRaw
			if result == expectedPartialFallback {
				// This is exactly a partial fallback, not a real context match
				return "", false
			}
		}

		return result, true
	}

	return "", false
}

// sortStrings sorts a slice of strings in place.
func sortStrings(strs []string) {
	sort.Strings(strs)
}

// buildV1Segment builds the v1 output segment for a node
func (t *Trie) buildV1Segment(node *TrieNode, part pathPart) string {
	if node.v1Name == "-" {
		// Skip this segment in output
		return ""
	}

	v1Name := node.v1Name
	if v1Name == "" {
		v1Name = part.name
	}

	// Append array index if present
	if part.arrayIndex != "" && node.isArray {
		return v1Name + part.arrayIndex
	}

	return v1Name
}

// buildRemainingPath constructs path from remaining parts
func (t *Trie) buildRemainingPath(parts []pathPart, startIndex int) string {
	var segments []string
	for i := startIndex; i < len(parts); i++ {
		segment := parts[i].name
		if parts[i].arrayIndex != "" {
			segment += parts[i].arrayIndex
		}
		segments = append(segments, segment)
	}
	return strings.Join(segments, ".")
}

// resolveStepType performs lazy step type resolution for context-aware matching.
// Resolution priority:
//  1. If convContext.StepType is already set, use it (pre-resolved)
//  2. If matchContext has a captured step ID, look it up in StepTypeMap
//  3. Fall back to convContext.CurrentStepType (the step we're inside)
//
// This enables efficient context resolution only when needed (at step.spec nodes).
func (t *Trie) resolveStepType(ctx *matchContext, convContext *ConversionContext) string {
	if convContext == nil {
		return ""
	}

	// Priority 1: Already resolved step type
	if convContext.StepType != "" {
		return convContext.StepType
	}

	// Priority 2: Look up captured step ID from path traversal
	if ctx.lastStepID != "" && convContext.StepTypeMap != nil {
		if stepType, ok := convContext.StepTypeMap[ctx.lastStepID]; ok {
			return stepType
		}
	}

	// Priority 3: Fall back to current step type (step we're inside)
	return convContext.CurrentStepType
}

// isStepInternalField checks if a field name indicates step-internal content
// that should trigger FQN building when UseFQN is enabled.
func isStepInternalField(fieldName string) bool {
	switch fieldName {
	case "spec", "output", "identifier", "name", "type", "timeout", "failureStrategies",
		"strategy", "when", "delegateSelectors", "description", "status":
		return true
	default:
		return false
	}
}
