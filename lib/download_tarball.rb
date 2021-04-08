# frozen_string_literal: true

require 'fileutils'
require 'json'
require 'time'

require_relative './common'
require_relative './commit'
require_relative './output'
require_relative './tarball'

class DownloadTarball
  include Common
  include Output

  logger_name 'download-tarball'

  TARBALL_URL = 'https://github.com/%s/tarball/%s'
  COMMIT_URL = 'https://api.github.com/repos/%s/commits/%s'

  attr_reader :ref
  attr_reader :repo
  attr_reader :output
  attr_reader :log_level

  def initialize(ref:, repo:, output:, log_level:)
    @ref = ref
    @repo = repo
    @output = output
    @log_level = log_level

    err 'branch/tag/sha argument cannot be empty' if ref.nil? || ref.empty?
  end

  def perform
    FileUtils.mkdir_p(output)
    tarball = Tarball.new(file: target, commit: commit)

    if File.exist?(target)
      info "#{filename} already exists locally, attempting to use."
      return tarball
    end

    info 'Downloading tarball from GitHub. This could take a while, ' \
         'please be patient.'
    result = run_cmd('curl', '-L', url, '-o', target)
    err 'Download failed.' if !result || !File.exist?(target)

    tarball
  end

  def url
    @url ||= format(TARBALL_URL, repo, commit.sha)
  end

  def filename
    @filename ||= "#{repo.gsub(/[^\w]/, '-')}-#{commit.sha_short}.tgz"
  end

  def target
    @target ||= File.join(output, filename)
  end

  def commit
    return @commit if @commit

    info "Fetching info for git ref: #{ref}"
    url = format(COMMIT_URL, repo, ref)
    commit_json = http_get(url)
    err "Failed to get commit info about: #{ref}" if commit_json.nil?

    parsed = JSON.parse(commit_json)
    commit = Commit.new(
      sha: parsed&.dig('sha'),
      time: Time.parse(parsed&.dig('commit', 'committer', 'date'))
    )

    err 'Failed to get commit SHA' if commit.sha.nil? || commit.sha.empty?
    err 'Failed to get commit time' if commit.time.nil?

    @commit = commit
  end
end
