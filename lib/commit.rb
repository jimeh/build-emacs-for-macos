# frozen_string_literal: true

require 'json'
require 'time'
require 'yaml'

require_relative './errors'
require_relative './common'

class Commit
  include Common

  attr_reader :repo
  attr_reader :ref
  attr_reader :message
  attr_reader :sha
  attr_reader :time

  def initialize(sha:, time:, repo: nil, ref: nil, message: nil)
    @sha = sha
    @time = time
    @repo = repo
    @ref = ref
    @message = message
  end

  def sha_short
    sha[0..6]
  end

  def to_hash
    {
      'repo' => repo,
      'ref' => ref,
      'sha' => sha,
      'sha_short' => sha_short,
      'time' => time.utc,
      'timestamp' => time.utc.to_i,
      'message' => message
    }
  end

  def to_yaml
    to_hash.to_yaml
  end
end
