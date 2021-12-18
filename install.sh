#!/usr/bin/env bash
set -e

{

  tfvmDir="$HOME/.tfvm/bin"

  # Check CPU architecture and OS
  arch=$(uname -m)
  case $arch in 

    x86_64)
      arch="amd64"
      ;;
    
    i386 | i686)
      arch="386"
      ;;

    aarch64_be | aarch64 | armv8b | armv8l | arm64)
      arch="arm64"
      ;;
    
    *)
      echo "ERROR: Unsupported architecture, try installing manually if you believe this is incorrect."
      trap exit ERR
      ;;
  
  esac

  os=$(uname -s)
  case $os in

    Linux)
      os="linux"
      ;;

    Darwin)
      os="darwin"
      ;;

    *)
      echo "ERROR: Unsupported OS, try installing manually if you believe this is incorrect."
      trap exit ERR
      ;;

  esac

  # Download latest version
  url=$(curl -s https://api.github.com/repos/ethanhassett/tfvm/releases | grep browser_download_url | grep $os | grep $arch | head -n 1 | cut -d '"' -f 4)
  pkg=$(echo $url | sed 's/.*\///')
  echo "Downloading $pkg from GitHub..."
  wget -qP /tmp $url

  # Verify checksum
  checksumUrl="https://github.com/ethanhassett/tfvm/releases/latest/download/checksum.txt"
  checksum=$(wget -q -O - $checksumUrl | cat | grep $pkg | head -n1 | cut -d " " -f1 | xargs)

  case $os in

    linux)
      calcsum=$(sha256sum /tmp/$pkg | head -n1 | cut -d " " -f1 | xargs)
      ;;

    darwin)
      calcsum=$(shasum -a 256 /tmp/$pkg | head -n1 | cut -d " " -f1 | xargs)
      ;;

    *)
      echo "ERROR: Unsupported OS, try installing manually if you believe this is incorrect."
      trap exit ERR
      ;;

  esac

  if [[ $calcsum != $checksum ]]; then
    echo "ERROR: Could not verify checksum. Please manually verify and install."
    trap exit ERR
  fi

  # Extract to /usr/bin/tfvm
  if ! [[ -d $tfvmDir ]]; then
    mkdir -p $tfvmDir
  fi
  tar -xzf /tmp/$pkg -C $tfvmDir
  rm /tmp/$pkg

  # Determines profile file to add tfvm directory to PATH
  profile=""
  if [[ "${PROFILE}" = '/dev/null' ]]; then
    exit 0
  fi

  if [[ -n "${PROFILE}" ]] && [[ -f "${PROFILE}" ]]; then
    profile=$PROFILE
  fi

  if [[ "${SHELL}" == *"bash"* ]]; then
    if [[ -f "$HOME/.bashrc" ]]; then
      profile="$HOME/.bashrc"
    elif [[ -f "$HOME/.bash_profile" ]]; then
      profile="$HOME/.bash_profile"
    fi
  elif [[ "${SHELL}" == *"zsh"* ]]; then
    profile="$HOME/.zshrc"
  fi

  if ! [[ -z $profile ]]; then
    # Add export line to file if it isn't already there
    if grep -Fq 'export PATH="$PATH:$HOME/.tfvm/bin"' $profile; then
      echo "tfvm v$latestVersion was installed successfully!"
    else
      printf "\n# Add tfvm to PATH\n" >> $profile
      echo 'export PATH="$PATH:$HOME/.tfvm/bin"' >> $profile
      source $profile
      echo "tfvm v$latestVersion was installed successfully!"
    fi
  fi
  
}