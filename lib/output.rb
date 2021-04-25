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
end
