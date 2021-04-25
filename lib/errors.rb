# frozen_string_literal: true

def handle_error(err)
  warn "ERROR: #{err.message}"
  Process.exit 1
end

class Error < StandardError; end

class CommitNotFound < Error; end
class NoCommitSHA < Error; end
class NoCommitTime < Error; end
