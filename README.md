# Slack Orb for CircleCI [![CircleCI Build Status](https://circleci.com/gh/CircleCI-Public/slack-orb.svg?style=shield "CircleCI Build Status")](https://circleci.com/gh/CircleCI-Public/slack-orb) [![CircleCI Orb Version](https://badges.circleci.com/orbs/circleci/slack.svg)](https://circleci.com/orbs/registry/orb/circleci/slack) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/circleci-public/slack-orb/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

The Slack Orb for CircleCI connects your CI/CD pipelines to Slack for ChatOps. Fully customizable, platform agnostic, and easy to integrate into your new or existing CI/CD pipelines.

## Slack Orb v5+ Breaking Changes

The Slack Orb for CircleCI is now written in Go! This means the orb is now more easily maintained, and better supports a wide range of platforms, no longer requiring BASH, CURL, JQ, or any other shell tools or commands.

1. The orb now downloads the Slack Orb CLI binary to `.circleci/orbs/circleci/slack/$PLATFORM/$ARCH` in the working directory.
   1. This can be checksum verified
   2. This is cacheable if desired.
2. You can no longer use evaluated sub-shell commands (e.g. `$(date +%s)`).
   1. Instead, you can pre-populate environment variables in the `$BASH_ENV` file. [See more in our wiki]()
   2. The Slack Orb has a small number of "built in" environment variables for use in templates [See more in our wiki]()
3. On-hold job removed. Low usage prompted removal from this release for simplicity of the documentation. We will revisit this as feedback comes in.

## Usage

### Setup

In order to use the Slack Orb on CircleCI you will need to create a Slack App and provide an OAuth token. 

**Full guide in the wiki:** [How to setup Slack orb](https://github.com/CircleCI-Public/slack-orb/wiki/Setup)

### Use In Config

For config examples, see the [Orb Registry listing](http://circleci.com/orbs/registry/orb/circleci/slack).

For more detailed information on customization, please see the [Slack Orb Wiki](https://github.com/CircleCI-Public/slack-orb/wiki).

## Templates

The Slack Orb comes with a number of included templates to get your started with minimal setup. Feel free to use an included template or create your own.

| Template Preview  | Template  | Description |
| ------------- | ------------- | ------------- |
| ![basic_success_1](./.github/img/basic_success_1.png)  | basic_success_1   | Should be used with the "pass" event. |
| ![basic_fail_1](./.github/img/basic_fail_1.png)  | basic_fail_1   | Should be used with the "fail" event. |
| ![success_tagged_deploy_1](./.github/img/success_tagged_deploy_1.png)  | success_tagged_deploy_1   | To be used in the event of a successful deployment job. _see orb [usage examples](https://circleci.com/developer/orbs/orb/circleci/slack#usage-examples)_ |


## Custom Message Template

  1. Open the Slack Block Kit Builder: https://app.slack.com/block-kit-builder/
  2. Design your desired notification message.
  3. Replace any placeholder values with $ENV environment variable strings.
  4. Set the resulting code as the value for your `custom` parameter.

  ```yaml
- slack/notify:
      event: always
      custom: |
        {
          "blocks": [
            {
              "type": "section",
              "fields": [
                {
                  "type": "plain_text",
                  "text": "*Notification from $CIRCLE_JOB*",
                  "emoji": true
                }
              ]
            }
          ]
        }
  ```


## FAQ

View the [FAQ in the wiki](https://github.com/CircleCI-Public/slack-orb/wiki/FAQ)

## Contributing

We welcome [issues](https://github.com/CircleCI-Public/slack-orb/issues) to and [pull requests](https://github.com/CircleCI-Public/slack-orb/pulls) against this repository!

For further questions/comments about this or other orbs, visit [CircleCI's orbs discussion forum](https://discuss.circleci.com/c/orbs).

### Developer Setup

This repository is configured as a monorepo, containing the `orb` source code, which is a collection of `BASH` and CircleCI `YAML`, and the `cli` source code, which is the main `go` binary utilized by this orb.

#### Go

1. Clone this repository
2. Install [taskfile.dev](https://taskfile.dev/installation/) if you do not have it.
   1. HOMEBREW: `brew install go-task/tap/go-task`
   2. NPM: `npm install -g @go-task/cli`
   3. CHOCOLATEY: `choco install go-task`
3. Run `task sync` to download dependencies.
4. Before pushing your branch, ensure to run `task tidy`