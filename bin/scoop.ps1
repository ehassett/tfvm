param (
  [switch] $uninstall
)

$tfvmDir = "$env:USERPROFILE\.tfvm"
$tfvmBinDir = "$tfvmDir\bin"
$tfvmRegex = [Regex]::Escape($tfvmBinDir)
$userPath = [Environment]::GetEnvironmentVariable("Path", "User").TrimEnd(";")
$userPathArray = $userPath -split ';' | Where-Object {$_ -notMatch "^$tfvmRegex\\?"}

if ($uninstall) {
  $newUserPath = ($userPathArray) -join ";"
  Remove-Item $tfvmDir -Recurse
} else {
  $newUserPath = ($userPathArray + $tfvmBinDir) -join ";"
}

[Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
$env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [Environment]::GetEnvironmentVariable("Path", "User")