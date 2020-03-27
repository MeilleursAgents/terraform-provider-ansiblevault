#!/usr/bin/env bash

set -o nounset -o pipefail -o errexit

github_last_release() {
  if [[ "${#}" -ne 1 ]]; then
    printf "%bUsage: github_last_release owner/repo%b\n" "${RED}" "${RESET}"
    return 1
  fi

  local RED="\033[31m"
  local RESET="\033[0m"

  local HTTP_OUTPUT="http_output.txt"
  local CLIENT_ARGS=("curl" "-q" "-sSL" "--max-time" "30" "-o" "${HTTP_OUTPUT}" "-w" "%{http_code}")
  if [[ -n ${GITHUB_OAUTH_TOKEN:-} ]]; then
    CLIENT_ARGS+=("-H" "Authorization: token ${GITHUB_OAUTH_TOKEN}")
  fi

  local LATEST_RELEASE
  LATEST_RELEASE="$("${CLIENT_ARGS[@]}" "https://api.github.com/repos/${1}/releases/latest")"
  if [[ "${LATEST_RELEASE}" != "200" ]]; then
    printf "%bUnable to list latest release for %s%b\n" "${RED}" "${1}" "${RESET}"
    cat "${HTTP_OUTPUT}" && rm "${HTTP_OUTPUT}"
    return
  fi

  python -c "import json; print(json.load(open('${HTTP_OUTPUT}'))['tag_name'])"
  rm "${HTTP_OUTPUT}"
}

main() {
  local BLUE="\033[34m"
  local RESET="\033[0m"

  local PLUGIN_VERSION
  PLUGIN_VERSION="$(github_last_release MeilleursAgents/terraform-provider-ansiblevault)"

  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

  if [[ "${ARCH}" = "x86_64" ]]; then
    ARCH="amd64"
  fi

  local PLUGIN_DIR="${HOME}/.terraform.d/plugins/${OS}_${ARCH}/"

  printf "%bInstalling terraform-provider-ansiblevault version %s into %s%b\n" "${BLUE}" "${PLUGIN_VERSION}" "${PLUGIN_DIR}" "${RESET}"

  mkdir -p "${PLUGIN_DIR}"
  (
    cd "${PLUGIN_DIR}" || return
    curl -q -sSL -o "terraform-provider-ansiblevault_${PLUGIN_VERSION}" "https://github.com/MeilleursAgents/terraform-provider-ansiblevault/releases/download/${PLUGIN_VERSION}/terraform-provider-ansiblevault_${OS}_${ARCH}_${PLUGIN_VERSION}"
    chmod +x "terraform-provider-ansiblevault_${PLUGIN_VERSION}"
  )
}

main
