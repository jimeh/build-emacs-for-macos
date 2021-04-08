# frozen_string_literal: true

require 'yaml'

class Commit
  attr_reader :sha
  attr_reader :time

  def initialize(sha:, time:)
    @sha = sha
    @time = time
  end

  def sha_short
    sha[0..6]
  end

  def to_hash
    {
      'sha' => sha,
      'sha_short' => sha_short,
      'time' => time
    }
  end

  def to_yaml
    to_hash.to_yaml
  end
end
