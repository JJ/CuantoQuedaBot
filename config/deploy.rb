set :application, "CuantoQuedaBot"
set :repo_url, "https://github.com/JJ/CuantoQuedaBot.git"
token = ENV['BOT_TOKEN']
puts "token is #{token}"
user = "nitrous"

set :user, user
set :deploy_to, "/home/#{user}/app"
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :symlinks, { "pids" => "pids" }

role :app, "jjmerelo-8639.nitrouspro.com"

task :production do
  server "jjmerelo-8639.nitrouspro.com", user: "nitrous", port: 32769
end

after 'deploy:updated', 'go:build'

namespace :go do
  task :build do
    on roles(:app) do
      execute "export GOROOT=/usr/local/opt/go;export GOPATH=/home/nitrous/lib/Go;export GOBIN=$GOPATH/bin;cd #{release_path};/usr/local/opt/go/bin/go get"
      execute "export GOROOT=/usr/local/opt/go;export GOPATH=/home/nitrous/lib/Go;export GOBIN=$GOPATH/bin;cd #{release_path};export BOT_TOKEN=#{token};/usr/local/opt/go/bin/go run CuantoQuedaBot.go"

    end
  end
end

