/**
 * commitlint configuration file
 * Documentation
 * https://commitlint.js.org/#/reference-rules
 * https://cz-git.qbb.sh/guide/
 */

module.exports = {
  // Extends rules
  extends: ['@commitlint/config-conventional'],
  // Custom rules
  rules: {
    // Type enum, git commit type must be one of the following types
    'type-enum': [
      2,
      'always',
      [
        'feat', // A new feature
        'fix', // A bug fix
        'docs', // Documentation only changes
        'style', // Changes that do not affect the meaning of the code
        'refactor', // A code change that neither fixes a bug nor adds a feature
        'perf', // A code change that improves performance
        'test', // Adding missing tests or correcting existing tests
        'build', // Changes that affect the build system or external dependencies
        'ci', // Changes to our CI configuration files and scripts
        'revert', // Reverts a previous commit
        'chore', // Other changes that don't modify src or test files
        'wip' // Work in progress
      ]
    ],
    'subject-case': [0], // No validation for subject case
    'type-case': [0], // No validation for type case
    'type-empty': [0], // Allow empty type
    'subject-empty': [0] // Allow empty subject
  },
  // Custom parser preset, support emoji
  parserPreset: {
    parserOpts: {
      headerPattern:
        /^(?:([\u00a9\u00ae\u2000-\u3300\ud83c\ud000-\udfff\ud83d\ud000-\udfff\ud83e\ud000-\udfff])\s)?([\w-]+)(?:\(([\w-]+)\))?:\s(.+)$/,
      headerCorrespondence: ['emoji', 'type', 'scope', 'subject']
    }
  },

  prompt: {
    messages: {
      type: 'Select the type of change that you\'re committing:',
      scope: 'Denote the SCOPE of this change (optional):',
      customScope: 'Denote the SCOPE of this change:',
      subject: 'Write a SHORT, IMPERATIVE tense description of the change:\n',
      body: 'Provide a LONGER description of the change (optional). Use "|" to break new line:\n',
      breaking: 'List any BREAKING CHANGES (optional). Use "|" to break new line:\n',
      footerPrefixesSelect: 'Select the ISSUES type of changeList by this change (optional):',
      customFooterPrefix: 'Input ISSUES prefix:',
      footer: 'List any ISSUES by this change. E.g.: #31, #34:\n',
      confirmCommit: 'Are you sure you want to proceed with the commit above?'
    },
    // prettier-ignore
    types: [
      { value: "✨ feat",     name: "feat:     ✨ A new feature", emoji: "✨" },
      { value: "🐛 fix",      name: "fix:      🐛 A bug fix", emoji: "🐛" },
      { value: "📝 docs",     name: "docs:     📝 Documentation only changes", emoji: "📝" },
      { value: "💄 style",    name: "style:    💄 Changes that do not affect the meaning of the code", emoji: "💄" },
      { value: "♻️ refactor", name: "refactor: ♻️ A code change that neither fixes a bug nor adds a feature", emoji: "♻️" },
      { value: "⚡️ perf",     name: "perf:     ⚡️ A code change that improves performance", emoji: "⚡️" },
      { value: "✅ test",     name: "test:     ✅ Adding missing tests or correcting existing tests", emoji: "✅" },
      { value: "📦️ build",    name: "build:    📦️ Changes that affect the build system or external dependencies", emoji: "📦️" },
      { value: "🎡 ci",       name: "ci:       🎡 Changes to our CI configuration files and scripts", emoji: "🎡" },
      { value: "⏪️ revert",   name: "revert:   ⏪️ Reverts a previous commit", emoji: "⏪️" },
      { value: "🔨 chore",    name: "chore:    🔨 Other changes that don't modify src or test files", emoji: "🔨" },
      { value: "🚧 wip",      name: "wip:      🚧 Work in progress", emoji: "🚧" },
    ],
    useEmoji: false,
    emojiAlign: 'center',
    themeColorCode: '',
    scopes: [
      { value: 'web', name: 'web: Frontend related' },
      { value: 'api', name: 'api: API interface' },
      { value: 'task', name: 'task: Task scheduling' },
      { value: 'node', name: 'node: Node management' },
      { value: 'auth', name: 'auth: Authentication' },
      { value: 'db', name: 'db: Database' },
      { value: 'config', name: 'config: Configuration' },
      { value: 'deps', name: 'deps: Dependencies update' }
    ],
    allowCustomScopes: true,
    allowEmptyScopes: true,
    customScopesAlign: 'bottom',
    customScopesAlias: 'custom',
    emptyScopesAlias: 'empty',
    upperCaseSubject: false,
    markBreakingChangeMode: false,
    allowBreakingChanges: ['feat', 'fix'],
    breaklineNumber: 100,
    breaklineChar: '|',
    skipQuestions: ['breaking', 'footerPrefix', 'footer'], // Skip these steps
    issuePrefixes: [{ value: 'closed', name: 'closed:   ISSUES has been processed' }],
    customIssuePrefixAlign: 'top',
    emptyIssuePrefixAlias: 'skip',
    customIssuePrefixAlias: 'custom',
    allowCustomIssuePrefix: true,
    allowEmptyIssuePrefix: true,
    confirmColorize: true,
    maxHeaderLength: Infinity,
    maxSubjectLength: Infinity,
    minSubjectLength: 0,
    scopeOverrides: undefined,
    defaultBody: '',
    defaultIssues: '',
    defaultScope: '',
    defaultSubject: ''
  }
}
