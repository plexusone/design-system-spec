package dss

// Governance defines versioning, contribution, and deprecation policies.
type Governance struct {
	// Versioning describes the versioning strategy.
	Versioning *VersioningPolicy `json:"versioning,omitempty"`

	// Contribution describes how to contribute to the design system.
	Contribution *ContributionPolicy `json:"contribution,omitempty"`

	// Deprecation describes how items are deprecated and removed.
	Deprecation *DeprecationPolicy `json:"deprecation,omitempty"`

	// DecisionProcess describes how design decisions are made.
	DecisionProcess *DecisionProcess `json:"decisionProcess,omitempty"`

	// Ownership describes who owns different parts of the system.
	Ownership []OwnershipRecord `json:"ownership,omitempty"`

	// SupportPolicy describes support commitments.
	SupportPolicy *SupportPolicy `json:"supportPolicy,omitempty"`
}

// VersioningPolicy describes the versioning strategy.
type VersioningPolicy struct {
	// Strategy is the versioning approach (semver, calver, etc.).
	Strategy string `json:"strategy"`

	// Description explains the versioning policy.
	Description string `json:"description,omitempty"`

	// BreakingChangePolicy explains when breaking changes are allowed.
	BreakingChangePolicy string `json:"breakingChangePolicy,omitempty"`

	// ReleaseSchedule describes the release cadence.
	ReleaseSchedule string `json:"releaseSchedule,omitempty"`

	// ChangelogLocation is the URL or path to the changelog.
	ChangelogLocation string `json:"changelogLocation,omitempty"`
}

// ContributionPolicy describes contribution guidelines.
type ContributionPolicy struct {
	// Description provides an overview of the contribution process.
	Description string `json:"description,omitempty"`

	// Guidelines lists contribution guidelines.
	Guidelines []string `json:"guidelines,omitempty"`

	// Workflow describes the contribution workflow.
	Workflow *ContributionWorkflow `json:"workflow,omitempty"`

	// Templates provides PR/issue templates.
	Templates []ContributionTemplate `json:"templates,omitempty"`

	// CodeOfConductURL links to the code of conduct.
	CodeOfConductURL string `json:"codeOfConductUrl,omitempty" jsonschema:"format=uri"`
}

// ContributionWorkflow describes the steps for contributing.
type ContributionWorkflow struct {
	// Steps lists the contribution steps.
	Steps []WorkflowStep `json:"steps"`

	// RequiredReviewers lists required reviewer roles.
	RequiredReviewers []string `json:"requiredReviewers,omitempty"`

	// AutomatedChecks lists automated checks that must pass.
	AutomatedChecks []string `json:"automatedChecks,omitempty"`
}

// WorkflowStep describes a step in the contribution workflow.
type WorkflowStep struct {
	// Order is the step number.
	Order int `json:"order"`

	// Name is the step name.
	Name string `json:"name"`

	// Description explains what happens in this step.
	Description string `json:"description"`

	// Assignee describes who is responsible for this step.
	Assignee string `json:"assignee,omitempty"`
}

// ContributionTemplate provides a template for contributions.
type ContributionTemplate struct {
	// Type is the template type (pr, issue, rfc).
	Type string `json:"type"`

	// Name is the template name.
	Name string `json:"name"`

	// URL links to the template file.
	URL string `json:"url,omitempty" jsonschema:"format=uri"`
}

// DeprecationPolicy describes how items are deprecated.
type DeprecationPolicy struct {
	// WarningPeriod is the minimum time items remain deprecated before removal.
	WarningPeriod string `json:"warningPeriod"`

	// NotificationChannels lists how deprecations are communicated.
	NotificationChannels []string `json:"notificationChannels,omitempty"`

	// MigrationGuideRequired indicates whether migration guides are required.
	MigrationGuideRequired bool `json:"migrationGuideRequired,omitempty"`

	// DeprecatedItems lists currently deprecated items.
	DeprecatedItems []DeprecatedItem `json:"deprecatedItems,omitempty"`
}

// DeprecatedItem describes a deprecated item in the design system.
type DeprecatedItem struct {
	// Type is the item type (component, token, pattern).
	Type string `json:"type"`

	// ID is the item identifier.
	ID string `json:"id"`

	// DeprecatedSince is when the item was deprecated.
	DeprecatedSince string `json:"deprecatedSince"`

	// RemovalDate is when the item will be removed.
	RemovalDate string `json:"removalDate,omitempty"`

	// Replacement is the ID of the replacement item.
	Replacement string `json:"replacement,omitempty"`

	// MigrationGuideURL links to migration documentation.
	MigrationGuideURL string `json:"migrationGuideUrl,omitempty" jsonschema:"format=uri"`

	// Reason explains why the item was deprecated.
	Reason string `json:"reason,omitempty"`
}

// DecisionProcess describes how design decisions are made.
type DecisionProcess struct {
	// Description provides an overview of the decision process.
	Description string `json:"description,omitempty"`

	// RFCRequired indicates whether RFCs are required for changes.
	RFCRequired bool `json:"rfcRequired,omitempty"`

	// ApprovalRequired lists who must approve changes.
	ApprovalRequired []string `json:"approvalRequired,omitempty"`

	// DecisionRecordLocation is where decisions are documented.
	DecisionRecordLocation string `json:"decisionRecordLocation,omitempty"`
}

// OwnershipRecord describes ownership of a part of the design system.
type OwnershipRecord struct {
	// Area is the area owned (e.g., "foundations", "components/forms").
	Area string `json:"area"`

	// Team is the owning team.
	Team string `json:"team"`

	// ContactEmail is the contact email for the team.
	ContactEmail string `json:"contactEmail,omitempty" jsonschema:"format=email"`

	// SlackChannel is the team's Slack channel.
	SlackChannel string `json:"slackChannel,omitempty"`
}

// SupportPolicy describes support commitments.
type SupportPolicy struct {
	// Description provides an overview of support.
	Description string `json:"description,omitempty"`

	// SupportedVersions lists which versions are currently supported.
	SupportedVersions []SupportedVersion `json:"supportedVersions,omitempty"`

	// SupportChannels lists how to get support.
	SupportChannels []string `json:"supportChannels,omitempty"`

	// ResponseTimeTargets describes expected response times.
	ResponseTimeTargets map[string]string `json:"responseTimeTargets,omitempty"`
}

// SupportedVersion describes a supported version.
type SupportedVersion struct {
	// Version is the version number.
	Version string `json:"version"`

	// Status is the support status (active, maintenance, end-of-life).
	Status string `json:"status"`

	// EndOfLifeDate is when support ends.
	EndOfLifeDate string `json:"endOfLifeDate,omitempty"`
}
