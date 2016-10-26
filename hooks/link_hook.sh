#!/bin/bash

local_path=`pwd`
ln -s "$local_path/hooks/pre-commit" "$local_path/.git/hooks/pre-commit"

