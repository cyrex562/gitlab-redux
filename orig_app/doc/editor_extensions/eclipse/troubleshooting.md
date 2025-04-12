---
stage: Create
group: Editor Extensions
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
description: Connect and use GitLab Duo in Eclipse.
title: Eclipse troubleshooting
---

{{< details >}}

- Tier: Free, Premium, Ultimate
- Offering: GitLab.com, GitLab Self-Managed, GitLab Dedicated
- Status: Experiment

{{< /details >}}

{{< alert type="disclaimer" />}}

If the steps on this page don't solve your problem, check the
[list of open issues](https://gitlab.com/gitlab-org/editor-extensions/gitlab-eclipse-plugin/-/issues/?sort=created_date&state=opened&first_page_size=100)
in the Eclipse plugin's project. If an issue matches your problem, update the issue.
If no issues match your problem, [create a new issue](https://gitlab.com/gitlab-org/editor-extensions/gitlab-eclipse-plugin/-/issues/new).

## Review the Error Log

1. In the menu bar of your IDE, select **Window**.
1. Expand **Show View**, then select **Error Log**.
1. Search for errors referencing the `gitlab-eclipse-plugin` plugins.

## Locate the Workspace Log file

The Workspace log file, named `.log` is located in the directory `<your-eclipse-workspace>/.metadata`.

## Enable GitLab Language Server debug logs

To enable GitLab Language Server debug logs:

1. In your IDE, select **Eclipse > Settings**.
1. On the left sidebar, select **GitLab**.
1. In **Language Server Log Level**, enter `debug`.
1. Select **Apply and Close**.

The debug logs are available in the `language_server.log` file. To view this file, either:

- Go to the directory `/Users/<user>/eclipse/<eclipse-version>/Eclipse.app/Contents/MacOS/.gitlab_plugin`, replacing `<user>` and `<eclipse-version>` with the appropriate values.
- Open the [Error logs](#review-the-error-log). Search for the log `Language server logs saved to: <file>.` where `<file>` is the absolute path to the `language_server.log` file.

## Certificate errors

{{< alert type="warning" >}}

You may experience errors connecting to GitLab if you connect to GitLab through a proxy or using custom certificates.
[Support for HTTP proxies](https://gitlab.com/gitlab-org/editor-extensions/gitlab-eclipse-plugin/-/issues/35)
and [support for custom certificates](https://gitlab.com/gitlab-org/editor-extensions/gitlab-eclipse-plugin/-/issues/36)
are proposed for a future release.

{{< /alert >}}
