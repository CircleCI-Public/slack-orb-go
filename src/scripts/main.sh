#!/bin/sh

# Determine the http client to use
# Returns 1 if no HTTP client is found
determine_http_client() {
  if command -v curl >/dev/null 2>&1; then
    HTTP_CLIENT=curl
  elif command -v wget >/dev/null 2>&1; then
    HTTP_CLIENT=wget
  else
    return 1
  fi
}

# Download a binary file
# $1: The path to save the file to
# $2: The URL to download the file from
# $3: The HTTP client to use (curl or wget)
download_binary() {
  if [ "$3" = "curl" ]; then
    set -x
    curl --fail --retry 3 -L -o "$1" "$2"
    set +x
  elif [ "$3" = "wget" ]; then
    set -x
    wget --tries=3 --timeout=10 --quiet -O "$1" "$2"
    set +x
  else
    return 1
  fi
}

detect_os() {
  detected_platform="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$detected_platform" in
  linux*) PLATFORM=Linux ;;
  darwin*) PLATFORM=Darwin ;;
  msys* | cygwin*) PLATFORM=Windows ;;
  *) return 1 ;;
  esac
}

detect_arch() {
  detected_arch="$(uname -m)"
  case "$detected_arch" in
  x86_64 | amd64) ARCH=x86_64 ;;
  i386 | i486 | i586 | i686) ARCH=i386 ;;
  arm64 | aarch64) ARCH=arm64 ;;
  arm*) ARCH=arm ;;
  *) return 1 ;;
  esac
}

# Determine the latest version of a GitHub release.
# $1: The GitHub organization
# $2: The GitHub repository
# $3: The HTTP client to use (curl or wget)
determine_release_latest_version() {
  url="https://github.com/$1/$2/releases/latest"

  if [ "$3" = "curl" ]; then
    LATEST_VERSION="$(curl --fail --retry 3 -Ls -o /dev/null -w '%{url_effective}' "https://github.com/$1/$2/releases/latest" | sed 's:.*/::')"
  elif [ "$3" = "wget" ]; then
    effective_url="$(wget --tries=3 --max-redirect=1000 --server-response -O /dev/null "$url" 2>&1 | awk '/Location: /{print $2}' | tail -1)"
    LATEST_VERSION="$(printf '%s' "$effective_url" | sed 's:.*/::')"
  else
    printf '%s\n' "Invalid HTTP client specified."
    return 1
  fi
}

# Print a warning message
# $1: The warning message to print
print_warn() {
  yellow="\033[1;33m"
  normal="\033[0m"
  printf "${yellow}%s${normal}\n" "$1"
}

# Print an error message
# $1: The error message to print
print_error() {
  red="\033[0;31m"
  normal="\033[0m"
  printf "${red}%s${normal}\n" "$1"
}

print_warn "This is an experimental version of the Slack Orb in Go."
print_warn "Thank you for trying it out and please provide feedback to us at https://github.com/CircleCI-Public/slack-orb-go/issues"

if ! detect_os; then
  printf '%s\n' "Unsupported operating system: $(uname -s)."
  exit 1
fi
printf '%s\n' "Operating system: $PLATFORM."

if ! detect_arch; then
  printf '%s\n' "Unsupported architecture: $(uname -m)."
  exit 1
fi
printf '%s\n' "Architecture: $ARCH."

base_dir="$(printf "%s" "$CIRCLE_WORKING_DIRECTORY" | sed "s|~|$HOME|")"
orb_bin_dir="$base_dir/.circleci/orbs/circleci/slack/$PLATFORM/$ARCH"
repo_org="CircleCI-Public"
repo_name="slack-orb-go"
binary="$orb_bin_dir/$repo_name"


if [ ! -f "$binary" ]; then
  print_warn "No binary found in $binary, downloading from $repo_url."
  print_warn "Consider caching the binary for future runs."
  mkdir -p "$orb_bin_dir"
  if ! determine_http_client; then
    print_error "cURL or wget is required to download the Slack binary."
    print_error "Please install cURL or wget and try again."
    exit 1
  fi
  printf '%s\n' "HTTP client: $HTTP_CLIENT."

  if ! determine_release_latest_version "$repo_org" "$repo_name" "$HTTP_CLIENT"; then
    print_error "Failed to determine latest version."
    exit 1
  fi
  printf '%s\n' "Release's latest version: $LATEST_VERSION."

  # TODO: Make the version configurable via command parameter
  repo_url="https://github.com/$repo_org/$repo_name/releases/download/$LATEST_VERSION/${repo_name}_${PLATFORM}_${ARCH}"
  [ "$PLATFORM" = "Windows" ] && repo_url="$repo_url.exe"
  printf '%s\n' "Release URL: $repo_url."

  # TODO: Check the sha256 of the downloaded binary
  if ! download_binary "$binary" "$repo_url" "$HTTP_CLIENT"; then
    print_error "Failed to download $repo_name binary from GitHub."
    exit 1
  fi

  printf '%s\n' "Downloaded $repo_name binary to $orb_bin_dir"
else
  printf '%s\n' "Skipping binary download since it already exists at $binary."
fi

# Validate the SHA256 checksum if provided
if [ -n "$SLACK_PARAM_SHA256" ]; then
  SHA256_SUBST=$(circleci env subst "$SLACK_PARAM_SHA256")
  printf '%s\n' "Input SHA256 checksum: $SHA256_SUBST"
  printf '%s\n' "Verifying $binary binary SHA256 checksum..."
  [ ! "$(which sha256sum)" ] && printf '%s\n' "sha256sum is required to verify the SHA256 checksum." && exit 1
  if ! printf '%s  %s' "$SHA256_SUBST" "$binary" | sha256sum -c -; then
    print_error "Failed to verify $binary binary SHA256 checksum."
    exit 1
  fi
fi

printf '%s\n' "Making $binary binary executable..."
if ! chmod +x "$binary"; then
  print_error "Failed to make $binary binary executable."
  exit 1
fi

printf '%s\n' "Executing \"$binary\" binary..."
set -x
"$binary"
exit_code=$?
set +x
if [ $exit_code -ne 0 ]; then
  printf '%s\n' "Failed to execute $binary binary or it exited with a non-zero exit code."
fi

printf '%s\n' "Removing $binary binary..."
rm -rf "$binary"
exit $exit_code
