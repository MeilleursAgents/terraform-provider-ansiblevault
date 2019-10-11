#!/usr/bin/env bash

github_last_release() {
  if [[ "${#}" -ne 1 ]]; then
    echo "Usage: github_last_release owner/repo"
    return 1
  fi

  local RED="\033[31m"
  local RESET="\033[0m"

  local OUTPUT_TXT="output.txt"
  local CLIENT_ARGS=("curl" "-q" "-sS" "-o" "${OUTPUT_TXT}" "-w" "%{http_code}")

  local LATEST_RELEASE="$("${CLIENT_ARGS[@]}" "https://api.github.com/repos/${1}/releases/latest")"
  if [[ "${LATEST_RELEASE}" != "200" ]]; then
    echo -e "${RED}Unable to list latest release for ${1}${RESET}"
    cat "${OUTPUT_TXT}" && rm "${OUTPUT_TXT}"
    return
  fi

  python -c "import json; print(json.load(open('${OUTPUT_TXT}'))['tag_name'])"
  rm "${OUTPUT_TXT}"
}

main() {
  local BLUE="\033[34m"
  local RESET="\033[0m"

  local PLUGIN_VERSION="$(github_last_release MeilleursAgents/terraform-provider-ansiblevault)"

  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

  if [[ "${ARCH}" = "x86_64" ]]; then
    ARCH="amd64"
  fi

  local PLUGIN_DIR="${HOME}/.terraform.d/plugins/${OS}_${ARCH}/"

  echo -e "${BLUE}Installing terraform-provider-ansiblevault version ${PLUGIN_VERSION} into ${PLUGIN_DIR}${RESET}"

  mkdir -p "${PLUGIN_DIR}"
  pushd "${PLUGIN_DIR}" || return
  curl -q -sS -Lo "terraform-provider-ansiblevault_${PLUGIN_VERSION}" "https://github.com/MeilleursAgents/terraform-provider-ansiblevault/releases/download/${PLUGIN_VERSION}/terraform-provider-ansiblevault_${OS}_${ARCH}_${PLUGIN_VERSION}"
  chmod +x "terraform-provider-ansiblevault_${PLUGIN_VERSION}"
  popd || return
}

main
