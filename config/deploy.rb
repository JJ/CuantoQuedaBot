set :application, "CuantoQuedaBot"
set :repo_url, "https://github.com/JJ/CuantoQuedaBot.git"
token = ENV['BOT_TOKEN']
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

    end
  end
  task :start do
    on roles(:app) do
      execute "export GOROOT=/usr/local/opt/go;export GOPATH=/home/nitrous/lib/Go;export GOBIN=$GOPATH/bin;cd #{release_path};export BOT_TOKEN=#{token};nohup /usr/local/opt/go/bin/go run CuantoQuedaBot.go > /home/nitrous/msgs.log 2>&1 &"
    end
  end

  task :stop do
    on roles(:app) do
      execute "pkill CuantoQuedaBot"
    end
  end
end


