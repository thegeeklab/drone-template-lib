local PipelineTest(deps=[],) = {
  kind: 'pipeline',
  name: 'test',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [
    {
      name: 'staticcheck',
      image: 'golang:1.16',
      commands: [
        'go run honnef.co/go/tools/cmd/staticcheck ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        },
      ],
    },
    {
      name: 'lint',
      image: 'golang:1.16',
      commands: [
        'go run golang.org/x/lint/golint -set_exit_status ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        },
      ],
    },
    {
      name: 'vet',
      image: 'golang:1.16',
      commands: [
        'go vet ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        },
      ],
    },
    {
      name: 'test',
      image: 'golang:1.16',
      commands: [
        'go test -cover ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        },
      ],
    },
  ],
  volumes: [
    {
      name: 'gopath',
      temp: {},
    },
  ],
  depends_on: deps,
  trigger: {
    ref: ['refs/heads/main', 'refs/tags/**', 'refs/pull/**'],
  },
};

local PipelineDocs(deps=[],) = {
  kind: 'pipeline',
  name: 'docs',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [
    {
      name: 'markdownlint',
      image: 'thegeeklab/markdownlint-cli',
      commands: [
        "markdownlint 'README.md' 'CONTRIBUTING.md'",
      ],
    },
    {
      name: 'spellcheck',
      image: 'node:lts-alpine',
      commands: [
        'npm install -g spellchecker-cli',
        "spellchecker --files 'README.md' 'CONTRIBUTING.md' -d .dictionary -p spell indefinite-article syntax-urls --no-suggestions",
      ],
      environment: {
        FORCE_COLOR: true,
        NPM_CONFIG_LOGLEVEL: 'error',
      },
    },
  ],
  depends_on: deps,
  trigger: {
    ref: ['refs/heads/main', 'refs/tags/**', 'refs/pull/**'],
  },
};

local PipelineNotifications(deps=[],) = {
  kind: 'pipeline',
  name: 'notifications',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [
    {
      name: 'matrix',
      image: 'plugins/matrix',
      settings: {
        homeserver: { from_secret: 'matrix_homeserver' },
        roomid: { from_secret: 'matrix_roomid' },
        template: 'Status: **{{ build.status }}**<br/> Build: [{{ repo.Owner }}/{{ repo.Name }}]({{ build.link }}) ({{ build.branch }}) by {{ build.author }}<br/> Message: {{ build.message }}',
        username: { from_secret: 'matrix_username' },
        password: { from_secret: 'matrix_password' },
      },
      when: {
        status: ['success', 'failure'],
      },
    },
  ],
  depends_on: deps,
  trigger: {
    ref: ['refs/heads/main', 'refs/tags/**'],
    status: ['success', 'failure'],
  },
};

[
  PipelineTest(),
  PipelineDocs(deps=['test']),
  PipelineNotifications(deps=['docs']),
]
