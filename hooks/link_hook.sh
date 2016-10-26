#!/bin/bash

local_path=`pwd`
rm "$local_path/.git/hooks/pre-commit.sample"
ln -s "$local_path/hooks/pre-commit" "$local_path/.git/hooks/pre-commit"

