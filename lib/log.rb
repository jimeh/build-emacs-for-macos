# frozen_string_literal: true

class Log
  extend Forwardable

  attr_reader :name
  attr_reader :level

  def initialize(name, level = :info)
    @name = name
    @level = level
  end

  def_delegators :logger, :debug, :info, :warn, :error, :fatal, :unkonwn

  private

  def logger
    @logger ||= Logger.new($stderr).tap do |l|
      l.progname = name
      l.level = level
      l.formatter = formatter
    end
  end

  def formatter
    proc do |severity, _datetime, progname, msg|
      "==> [#{progname}] #{severity}: #{msg}\n"
    end
  end
end
