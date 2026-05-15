#compdef engine-cli
# Zsh completion script for Sovereign Engine CLI

_engine_cli() {
    local -a commands subcommands options

    commands=(
        'deploy: Deploy infrastructure to cloud provider'
        'cleanup: Clean up idle resources'
        'status: Check status of deployments'
        'notify: Send notifications'
        'budget: Manage budget limits'
        'compose: Compose infrastructure from YAML'
        'install: Install Sovereign Engine'
        'help: Show help information'
    )

    case $words[1] in
        deploy)
            subcommands=(
                '--provider:Cloud provider (aws, azure, gcp, hetzner, ovh, digitalocean)'
                '--region:Region for deployment'
                '--name:Deployment name'
                '--environment:Environment (dev, test, staging, production)'
                '--dry-run:Preview changes without applying'
            )
            ;;
        cleanup)
            subcommands=(
                '--environment:Environment to clean up'
                '--dry-run:Preview changes without applying'
                '--force:Force cleanup without checks'
                '--exclude-pattern:Pattern to exclude from cleanup'
            )
            ;;
        status)
            subcommands=(
                '--provider:Cloud provider'
                '--environment:Environment to check'
                '--json:Output in JSON format'
            )
            ;;
        notify)
            subcommands=(
                '--webhook:Webhook URL'
                '--slack:Slack channel'
                '--pagerduty:PagerDuty service'
                '--message:Custom message'
            )
            ;;
        budget)
            subcommands=(
                '--environment:Environment'
                '--limit:Budget limit'
                '--alert-threshold:Alert threshold'
                '--period:Time period'
            )
            ;;
        compose)
            subcommands=(
                '--file:YAML composition file'
                '--output:Output format'
                '--validate:Validate composition only'
                '--dry-run:Preview changes'
            )
            ;;
        install)
            subcommands=(
                '--version:Specific version to install'
                '--prefix:Installation prefix'
                '--path:Installation path'
            )
            ;;
        *)
            subcommands=()
            ;;
    esac

    if [[ $CURRENT_WORD == -* ]]; then
        _describe -t options 'options'
    elif [[ ${#words} -eq 1 ]]; then
        _describe -t commands 'commands'
    else
        _describe -t subcommands 'subcommands'
    fi
}

_engine_cli "$@"
