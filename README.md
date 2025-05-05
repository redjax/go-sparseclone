# Sparseclone

A Go module for automating my common flow for [git sparse checkouts](https://git-scm.com/docs/git-sparse-checkout).

When working with large repositories, a sparse checkout can help you 'scope' your git clone to avoid pulling in unnecessary code. If you have many subdirectories and only need to work on or use a small section of your repository, a sparse checkout can pull just that code from the remote repository.

This is an example of a sparse clone for my [Docker templates repository](https://github.com/redjax/docker_templates). In this example, I will pull my [`docker_gickup` container](https://github.com/redjax/docker_templates/tree/main/templates/backup/docker_gickup):

```shell
## Clone repo without checking out any code
git clone --no-checkout git@github.com:redjax/docker_templates.git docker_gickup

## Change directory into cloned code
cd docker_gickup

## Start the sparse checkout
git sparse-checkout init --cone

## Set the sparse clone to pull the script/ directory and the gickup template
git sparse-checkout set scripts templates/backup/docker_gickup

## Finally, checkout the feature branch where I'm working on gickup
#  and pull the code into the local repo
git checkout feat/gickup
```

## Usage

Run `sparseclone --help` to see all options.

Use the following flags to build your sparse checkout flow:

| Flag                                    | Description                                                                                                                                                                                                                                       |
| --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--provider [github, gitlab, codeberg]` | Set the remote provider where your repository is hosted.                                                                                                                                                                                          |
| `-u`/`--user`                           | Set your username on the remote, i.e. `https://github.com/<your-username>`                                                                                                                                                                        |
| `-r`/`--repo`                           | Set the name of your repository, i.e. `https://github.com/<your-username>/<your-repo-name>`                                                                                                                                                       |
| `-o`/`--output`                         | Set the directory where your code will be cloned. If no output is provided, the repository name will be used.                                                                                                                                     |
| `-b`/`--branch`                         | Set the branch to checkout after initializing the sparse clone.                                                                                                                                                                                   |
| `-p`/`--path`                           | Set 1 or more paths that exist in the repository to checkout with your sparse checkout. All files in the root of the repository will be cloned, but you can use `-p` to specify individual directories that you want present in your local clone. |
