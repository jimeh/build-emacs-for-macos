# frozen_string_literal: true

require 'yaml'

class Tarball
  attr_reader :file
  attr_reader :commit

  def initialize(file:, commit:)
    @file = file
    @commit = commit
  end

  def to_hash
    {
      'file' => file,
      'commit' => commit.to_hash
    }
  end

  def to_yaml
    to_hash.to_yaml
  end
end
