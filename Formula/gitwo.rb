class Gitwo < Formula
  desc "Git Worktree helper CLI with auto-cd and project detection"
  homepage "https://github.com/gitwohq/gitwo"
  version "0.1.0"
  
  # We'll need to set up releases with proper URLs
  # For now, this is a template that will need actual release URLs
  url "https://github.com/gitwohq/gitwo/releases/download/v0.1.0/gitwo_0.1.0_darwin_amd64.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  
  # Platform-specific binaries
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/gitwohq/gitwo/releases/download/v0.1.0/gitwo_0.1.0_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    else
      url "https://github.com/gitwohq/gitwo/releases/download/v0.1.0/gitwo_0.1.0_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end
  end
  
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/gitwohq/gitwo/releases/download/v0.1.0/gitwo_0.1.0_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/gitwohq/gitwo/releases/download/v0.1.0/gitwo_0.1.0_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "gitwo"
    
    # Install shell completions
    bash_completion.install "completions/gitwo.bash" => "gitwo"
    zsh_completion.install "completions/gitwo.zsh" => "_gitwo"
    fish_completion.install "completions/gitwo.fish" => "gitwo.fish"
  end

  def caveats
    <<~EOS
      🎉 gitwo has been installed!
      
      To enable auto-cd after 'gitwo new', install the shell wrapper:
        gitwo shell-install
      
      For current session only:
        eval "$(gitwo shell-init)"                   # bash/zsh
        gitwo shell-init --shell fish | source       # fish
        gitwo shell-init --shell pwsh | Out-String | Invoke-Expression  # PowerShell
      
      Quick start:
        gitwo init --with-shell                      # Initialize project with shell wrapper
        gitwo new my-feature                         # Create worktree and auto-cd
      
      Documentation: https://github.com/gitwohq/gitwo
    EOS
  end

  test do
    system "#{bin}/gitwo", "--version"
  end
end
