# General Development Rules

We are currently in architectural design phase using domain driven design. We are operating with a highest level domain definitions with DOMAIN.md files for each domain which contain the high level details of each domain. Each domain also has a series of supplemental files which contain specific low level details such as formulas, algorithms, data structures, interfaces, contracts, and other implementation specific details that we want to adhere to.

You should do task-based development. When defining tasks, identify what methods, classes, data structures, etc you will be modifying, then identify all the relevant files that you need to engage with. Before marking a task complete, you should verify that all references to the objects modified have been updated.

ALWAYS respect domain boundaries. When designing domain models, carefully think about whether the component you are adding belongs in the domain you are adding it to or if it should go in a different domain.

## Documentation Structure

### DOMAIN.md (Core Design)
- High-level domain concepts and boundaries
- Architectural diagrams and workflows
- Key interfaces and abstractions
- Integration patterns
- Strategic design decisions

*Code Examples Only For:*
- Critical interface definitions
- Core value objects
- Example integration patterns

### SUPPLEMENT-*.md (Reference)
- Data structures and schemas
- Algorithms and formulas
- Protocol/packet definitions
- State machines
- Reference data and constants

### INTERFACES.md
1. **Per-Subdomain Interfaces**
   - Each subdomain must have INTERFACES.md documenting:
     ```markdown
     ## Core Interfaces
     \`\`\`go
     type EntityRepository interface {
         Get(id EntityID) (Entity, error)
         Query(predicate func(Entity) bool) []Entity
     }
     \`\`\`
     
     ## Cross-Domain Contracts
     \`\`\`mermaid
     graph LR
         EM[EntityMgmt] -->|Provides| AI[AI System]
     \`\`\`
     ```

2. **Epic-Level Contracts**
   - Each domain epic (02-Entity etc) requires CONTRACTS.md:
     ```markdown
     ## Entity <> Spatial Indexing
     \`\`\`go
     type RangeQuery interface {
         EntitiesInRadius(pos Position, radius float64) []Entity
     }
     \`\`\`
     ```

3. **Interface Versioning**
   - Track stability in DOMAIN.md tables:
     | Interface          | Version | Status     | Consumers         |
     |--------------------|---------|------------|-------------------|
     | EntityRepository   | v1.2.3  | Frozen     | AI, Network, UI   |

4. **Cross-Validation**
   - Maintain INTERFACE_VALIDATION.md at root with:
     - Change approval process
     - Versioning rules
     - Contract testing requirements

## Workflow Rules

1. **Task-Based Development**
   - Identify all affected components before starting
   - Verify cross-references before marking complete
   - Use grep to find relevant files:
     `grep -r "pattern" ./<path>`

2. **Task Tracking**
   - `[ ]` = Todo
   - `[-]` = Partially complete not being actively worked on
   - `[O]` = Active Work in Progress
   - `[X]` = Completed

3. **Memory Maintenance**
   - Always update MEMORY.md each turn
   - Track state and open questions
   - When MEMORY.md gets too long, take a turn to prune and condense it
   - Never include code - only status

## Phase Rules

1. **Design Phase**
   - No functional code yet
   - Focus on interface contracts
   - Document Perl compatibility
   - DO NOT create any *.go files
   - When design has reached a point where we can move to implementation, tell the user

2. **Implementation**
   - Start with core interfaces
   - Add memory/supplement docs first
   - Verify against original behavior
   - DO NOT create any *.go files

3. **Integration**
   - Document cross-domain contracts
   - Update interface versions
   - Add compatibility tests
   - DO NOT create any *.go files

## Commit Rules

1. When a design is complete:
   - Update all TODO markers
   - Verify MEMORY.md status
   - Write descriptive commit message
   - Request human review

2. Never commit:
   - Unreviewed interface changes
   - Undocumented behavior
   - Breaking changes without versioning
