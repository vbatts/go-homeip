#!/usr/bin/env ruby 

require 'gdbm'

require 'rubygems'
require 'sinatra'

$DB = GDBM.new('/tmp/ips.db')

get '/' do
  'It works!'
end

get '/hosts' do
  $DB.keys.join("<br/>\n")
end

get '/ip/:host' do
  if $DB.has_key?(params[:host])
    '%s|%s' % [params[:host], $DB[params[:host]]]
  else
    status 218
  end
end

put '/ip/:host' do
  $DB[ params[:host] ] = request.ip
end

