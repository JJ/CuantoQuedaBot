set :application, "CuantoQuedaBot"
set :repository, "https://github.com/JJ/CuantoQuedaBot.git"

set :user, 'nitrous'
set :deploy_to, "/home/#{:user}/app"
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :symlinks, { "pids" => "pids" }


task :production do
  server "jjmerelo-8639.nitrouspro.com", user: "nitrous", port: 32769
end

after 'deploy:updated', 'go:build'

namespace :go do
  task :build do
    with_env('GOPATH', gopath) do
      run "go get -u github.com/JJ/CuantoQuedaBot"
      run "mkdir #{release_path}/bin"
      run "cp /home/#{:user}/go/bin/CuantoQuedaBot #{release_path}/bin/"
    end
  end
end

