# frozen_string_literal: true

def handle_error(err)
  warn "ERROR: #{err.message}"
  Process.exit 1
end

class Error < StandardError; end
