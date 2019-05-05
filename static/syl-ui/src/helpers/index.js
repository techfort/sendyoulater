export default promise => promise.then(data => ({ error: null, data }))
  .catch(error => ({ error, data: null }));