#!/bin/bash
input="$1"
{ grep '^#' "$input"; grep -v '^#' "$input" | sort -V; } > "$input".new
