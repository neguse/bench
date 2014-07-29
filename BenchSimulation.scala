package bench

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._
import io.gatling.http.Headers.Names._
import io.gatling.http.request._
import scala.concurrent.duration._
import bootstrap._
import assertions._

class BenchSimulation extends Simulation {

	val httpProtocol = http
		.baseURL("https://10.0.0.125")
		.acceptCharsetHeader("ISO-8859-1,utf-8;q=0.7,*;q=0.7")
		.acceptHeader("text/html,application/xhtml+xml,application/xml,application/json;q=0.9,*/*;q=0.8")
		.acceptEncodingHeader("gzip, deflate")
		.acceptLanguageHeader("fr,fr-fr;q=0.8,en-us;q=0.5,en;q=0.3")
		.disableFollowRedirect

	val headers = Map(
		"Connection" -> "Close",
    "Content-Type" -> "application/json")

	val scn = scenario("Scenario name")
    .during(100 seconds) {
      exec(
        http("index")
          .post("/jsonrpc")
          .body(new StringBody("""{"jsonrpc":"2.0", "id":1,"method":"index","params":[1,2,3]}"""))
          .headers(headers)
          .check(status.is(200)))
      .pause(10 seconds, 20 seconds)
    }

	setUp(scn.inject(ramp(10000 users) over (10 seconds)))
		.protocols(httpProtocol)
}
