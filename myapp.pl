#!/usr/bin/env perl
use Mojolicious::Lite;

use MojoX::JSON::RPC::Service;

my $svc = MojoX::JSON::RPC::Service->new;
$svc->register(
  'sum' => sub {
      my @params = @_;
      my $sum = 0;
      $sum += $_ for @params;
      return $sum;
  }
);
$svc->register(
  'index' => sub {
      my @params = @_;
      return "ok";
  }
);

plugin 'json_rpc_dispatcher' => {
  services => { 'jsonrpc' => $svc, },
};

app->start;
__DATA__

@@ index.html.ep
% layout 'default';
% title 'Welcome';
Welcome to the Mojolicious real-time web framework!

@@ layouts/default.html.ep
<!DOCTYPE html>
<html>
  <head><title><%= title %></title></head>
  <body><%= content %></body>
</html>
