$builtAt = Get-Date -Format "yyyy-MM-dd HH:mm zz00"
$gitAuthor = "vscodev"

if ($args.Count -gt 0 -and $args[0] -match "^v[0-9]+\.[0-9]+\.[0-9]+$") {
    $version = $args[0].TrimStart('v')
} else {
    $version = "dev"
}

$ldflags = @(
    "-s",
    "-w",
    "-X 'github.com/vscodev/alist/v3/internal/conf.BuiltAt=$builtAt'",
    "-X 'github.com/vscodev/alist/v3/internal/conf.GitAuthor=$gitAuthor'",
    "-X 'github.com/vscodev/alist/v3/internal/conf.Version=$version'"
) -join " "

function Build {
    param()
    $env:CGO_ENABLED = "1"
    go build -ldflags="$ldflags"
}

Build
