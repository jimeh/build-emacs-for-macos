# frozen_string_literal: true

require 'net/http'

module Common
  private

  def run_cmd(*args)
    info "executing: #{args.join(' ')}"
    system(*args) || err("Exit code: #{$CHILD_STATUS.exitstatus}")
  end

  def http_get(url)
    response = Net::HTTP.get_response(URI.parse(url))
    return unless response.code == '200'

    response.body
  end
end
