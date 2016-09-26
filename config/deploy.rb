set :application, "CuantoQuedaBot"
set :repo_url, "https://github.com/JJ/CuantoQuedaBot.git"

set :user, 'nitrous'
set :deploy_to, "/home/#{:user}/app"
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :symlinks, { "pids" => "pids" }

role :app, "jjmerelo-8639.nitrouspro.com"

task :production do
  server "jjmerelo-8639.nitrouspro.com", user: "nitrous", port: 32769
end

after 'deploy:updated', 'go:build'

set :default_env, { GOPATH: "/home/nitrous/lib/go" }

namespace :go do
  task :build do
    on roles(:app) do 
      execute "go get -u github.com/JJ/CuantoQuedaBot"
      execute "mkdir #{release_path}/bin"
      execute "cp /home/#{:user}/go/bin/CuantoQuedaBot #{release_path}/bin/"
    end
  end
end

