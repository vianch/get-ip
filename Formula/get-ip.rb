class GetIp < Formula
  desc "TUI for inspecting network interfaces, IPs, MACs, gateways, and DHCP mode"
  homepage "https://github.com/vianch/get-ip"
  url "https://github.com/vianch/get-ip/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "9c4e4b0d9317f0a3d04efab868bfdef02c8a81b21e520b3ff8cb66deada84cf8"
  license "MIT"
  head "https://github.com/vianch/get-ip.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X main.version=#{version}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags), "./"
  end

  test do
    assert_match "get-ip", shell_output("#{bin}/get-ip --version")
  end
end
