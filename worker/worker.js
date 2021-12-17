const { run, makeWorkerUtils } = require("graphile-worker");

run({
  taskDirectory: `${__dirname}/tasks`,
  concurrency: 2,
}).then(runner => {
  runner.events.on("job:success", ({ worker, job }) => {
    console.log(`--- Worker ${worker.workerId} completed job ${job.id}`);
  });
})
// 
// makeWorkerUtils({
// }).then(workerUtils => {
//   workerUtils.addJob('hello', { name: 'Test Worker is Running' })
// })
