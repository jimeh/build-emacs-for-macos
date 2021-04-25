# frozen_string_literal: true

require 'fileutils'
require 'json'
require 'time'

require_relative './base_action'

require_relative './commit_info'
require_relative './tarball'

class DownloadTarball < BaseAction
  TARBALL_URL = 'https://github.com/%s/tarball/%s'

  attr_reader :ref
  attr_reader :repo
  attr_reader :output
  attr_reader :logger

  def initialize(ref:, repo:, output:, logger:)
    @ref = ref
    @repo = repo
    @output = output
    @logger = logger

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
    @commit ||= CommitInfo.new(ref: ref, repo: repo, logger: logger).perform
  end
end
