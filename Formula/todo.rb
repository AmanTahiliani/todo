class Todo < Formula
  desc "A simple CLI todo application written in Go"
  homepage "https://github.com/AmanTahiliani/todo"
  version "0.1.0"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"todo", "./cmd/todo"
  end

  test do
    system "#{bin}/todo", "--help"
  end
end