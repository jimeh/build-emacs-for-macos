# frozen_string_literal: true

require 'forwardable'
require 'logger'

require_relative './errors'

module Output
  extend Forwardable

  def self.included(base)
    base.extend(ClassMethods)
  end

  module ClassMethods
    def logger_name(name = nil)
      return @logger_name if name.nil?

      @logger_name = name
    end
  end

  def_delegators :logger, :debug, :info, :warn, :error, :fatal, :unkonwn

  def err(msg = nil)
    raise Error, msg
  end

  private

  # override to set custom log level
  def log_level
    :info
  end

  def logger
    @logger ||= Logger.new($stderr).tap do |l|
      l.progname = self.class.logger_name
      l.level = log_level
      l.formatter = log_formatter
    end
  end

  def log_formatter
    proc do |severity, _datetime, progname, msg|
      "==> [#{progname}] #{severity}: #{msg}\n"
    end
  end
end
