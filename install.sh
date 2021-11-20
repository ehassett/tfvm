#!/usr/bin/env bash
set -e

{

  tfvmDir="$HOME/.tfvm/bin"

  # Get latest available verison from github
  IN=$(wget -q -O - https://raw.githubusercontent.com/ethanhassett/tfvm/main/tfvm.json | grep "\"version\":" | xargs)
  IFS=':' read -ra ADDR <<< $IN
  latestVersion=$(echo ${ADDR[1]} | tr -d ,)

  # Check OS
  os=$(uname -s)
  case $os in

    Darwin)
      os="darwin"
      ;;

    Linux)
      os="linux"
      ;;

    *)
      echo "ERROR: Unsupported oeprating system, try installing manually if you believe this is incorrect."
      trap exit ERR
      ;;

  esac

  # Check CPU architecture
  arch=$(uname -m)
  case $arch in 

    x86_64)
      arch="amd64"
      ;;
    
    i386 | i686)
      arch="386"
      ;;

    arm)
      arch="arm"
      ;;

    aarch64_be | aarch64 | armv8b | armv8l | arm64)
      arch="arm64"
      ;;
    
    *)
      echo "ERROR: Unsupported architecture, try installing manually if you believe this is incorrect."
      trap exit ERR
      ;;
  
  esac

  # Download latest version
  url="https://github.com/ethanhassett/tfvm/releases/download/$latestVersion/tfvm-$latestVersion-$os-$arch.tar.gz"
  checksumUrl="$url.md5"

  echo "Downloading from $url..."
  wget -q -O /tmp/tfvm-$latestVersion-$os-$arch.tar.gz $url

  # Verify checksum
  checksum=$(wget -q -O - $checksumUrl | cat)
  md5sum -c <<<"$checksum /tmp/tfvm-$latestVersion-$os-$arch.tar.gz"

  # Extract to /usr/bin/tfvm
  if ! [[ -d $tfvmDir ]]; then
    mkdir -p $tfvmDir
  fi
  tar -xzf /tmp/tfvm-$latestVersion-$os-$arch.tar.gz -C $tfvmDir
  rm /tmp/tfvm-$latestVersion-$os-$arch.tar.gz

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
