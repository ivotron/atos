# vio

Versioning for input/output files.

When working with a version-controlled project, we often use/obtain 
artifacts (conf files, logs, figures, etc.) for/from programs that 
correspond to a particular version of the project. After a couple of 
executions, it quickly becomes difficult to keep track of what 
versions of the project consumed/generated which files. `vio` helps to 
deal with this issue by allowing a user to create a snapshot of the 
unversioned files after a program has executed, and to store and 
associate this snapshot with the latest revision of the project.

## Example

```bash
git clone https://project.git

cd project

# work, work, work
git add -u
git commit -m "I worked hard and implemented many things"

# execute and generate some results
exec program -c params.conf

# commit anything that is not being tracked by git
vio commit -m "the result of my hard work"
```

## High-level

In a nutshell, vio:

 1. Finds all files that are not tracked by the VCS.
 2. Creates a dataset of all unversioned files.
 3. Puts the dataset in a storage backend, associating it to an 
    execution ID (`commit_id + timestamp`).
 4. Provides versioning-semantics for datasets, allowing users to 
    compare between distinct versions.
 5. Stores metadata for a dataset, allowing users to annotate and 
    contextualize unversioned files for future introspection.

The vio's "database" has the following schema:

```
 commit_id | execution_id | vio_commit_message | files | metadata |
```

## Multiple executions

One common use case is to compare results from multiple executions:

```
vio log --pretty=oneline

ca82a6df:20151123:120354 results with some conf1
ca82a6df:20151123:184832 and now with conf2
```

**TODO**

# vio vs. other tools

## `git-lfs`

`git-lfs` allows the inclusion of large files into a git repo. The 
main difference between vio and `git-lfs` is that `vio` lets you 
associate multiple datasets (or filesystem snapshots) to a single 
version of the git repo, while `git-lfs` can only associate a single 
one. In other words, the relationship between git commits and commits 
in the storage backend is one-to-one for `git-lfs` while one-to-many 
for `vio`.

Given the above, `vio` can use `git-lfs` as a backend, in the same way 
that the `git` backend is used by `vio`.

Other tools such as `git-annex`, etc. also fall in this category.

## artifact repositories

**TODO**

## CI tools

**TODO**


# references

Some use cases that this tool is aimed at solving:
  * <http://stackoverflow.com/q/18734739>
  * <http://academia.stackexchange.com/q/8359>
  * <http://academia.stackexchange.com/q/36995>
