# frozen_string_literal: true

require_relative './base_action'
require_relative './commit'

class CommitInfo < BaseAction
  COMMIT_URL = 'https://api.github.com/repos/%s/commits/%s'

  attr_reader :ref
  attr_reader :repo
  attr_reader :logger

  def initialize(ref:, repo:, logger:)
    @ref = ref
    @repo = repo
    @logger = logger

    err 'branch/tag/sha argument cannot be empty' if ref.nil? || ref.empty?
  end

  def perform
    info "Fetching info for git ref: #{ref}"

    url = format(COMMIT_URL, repo, ref)
    commit_json = http_get(url)

    err "Failed to get commit info about: #{ref}" if commit_json.nil?

    parsed = JSON.parse(commit_json)
    commit = Commit.new(
      repo: repo,
      ref: ref,
      sha: parsed&.dig('sha'),
      message: parsed&.dig('commit', 'message'),
      time: Time.parse(parsed&.dig('commit', 'committer', 'date')).utc
    )

    err 'Failed to get commit SHA' if commit.sha.nil? || commit.sha.empty?
    err 'Failed to get commit time' if commit.time.nil?

    commit
  end
end
