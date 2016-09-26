set :application, "CuantoQuedaBot"
set :repository, "https://github.com/JJ/CuantoQuedaBot.git"

set :user, 'nitrous'
set :deploy_to, "/home/#{:user}/app"
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :symlinks, { "pids" => "pids" }

namespace :go do
  task :production do
    server "jjmerelo-8639.nitrouspro.com", user: "nitrous", port: 32769
  end
end

