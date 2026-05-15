#!/bin/bash
# Bash completion script for Sovereign Engine CLI

_engine_cli_completion() {
    local cur prev words cword
    _init_completion || return

    case $prev in
        engine-cli)
            COMPREPLY=($(compgen -W "deploy cleanup status notify budget compose install help" -- "${cur}"))
            ;;
        deploy)
            COMPREPLY=($(compgen -W "--provider --region --name --environment --dry-run" -- "${cur}"))
            ;;
        cleanup)
            COMPREPLY=($(compgen -W "--environment --dry-run --force --exclude-pattern" -- "${cur}"))
            ;;
        status)
            COMPREPLY=($(compgen -W "--provider --environment --json" -- "${cur}"))
            ;;
        notify)
            COMPREPLY=($(compgen -W "--webhook --slack --pagerduty --message" -- "${cur}"))
            ;;
        budget)
            COMPREPLY=($(compgen -W "--environment --limit --alert-threshold --period" -- "${cur}"))
            ;;
        compose)
            COMPREPLY=($(compgen -W "--file --output --validate --dry-run" -- "${cur}"))
            ;;
        install)
            COMPREPLY=($(compgen -W "--version --prefix --path" -- "${cur}"))
            ;;
        --provider)
            COMPREPLY=($(compgen -W "aws azure gcp hetzner ovh digitalocean" -- "${cur}"))
            ;;
        --environment)
            COMPREPLY=($(compgen -W "dev test staging production" -- "${cur}"))
            ;;
        --region)
            COMPREPLY=($(compgen -W "us-east-1 us-west-2 eu-west-1 ap-southeast-1" -- "${cur}"))
            ;;
        *)
            COMPREPLY=($(compgen -f -- "${cur}"))
            ;;
    esac
}

complete -F _engine_cli_completion engine-cli
