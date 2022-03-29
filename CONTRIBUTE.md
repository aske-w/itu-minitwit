## Which repository setup will we use?

Github repository initialized as a non-bare repository from the minitwit files from the first lecture.

## Which branching model will we use?

We will use the Git Flow branching model. We will have a main branch representing the production code and a development branch from where developers branch out of and merge into, and from where the main branch will pull from. Further, all feature branches will be prepended with "feature/..." and hotfixes with "hotfix/...". We are not going to use release branches. We will instead create releases through Github.

## Which distributed development workflow will we use?

We want to use a combination of the Centralized Workflow and Integration-Manager workflow.[^1] This is in order to match our choice of the Git flow branching model. We are going to have two centralized repositories, the main branch representing the production code, and a development branch from where developers will branch out of. In contrast to the Integration-Manager Workflow we do not want each developer to create their own forks of the original repository. Instead everyone will work out of the original repository.

[^1]: https://git-scm.com/book/en/v2/Distributed-Git-Distributed-Workflows

## How do we expect contributions to look like?

The contributions should include a short and concise summary and a discription - if the contribution needs to be clarified or explained. 

## Who is responsible for integrating/reviewing contributions?

One or two people should review pull requests on the development branch such that the development branch is ready to deploy. We will set up pipelines which will test and build any updates to the production code base, and reject them if the tests or build process fails.
