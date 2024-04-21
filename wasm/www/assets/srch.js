// Get srch Options
async function getCfg() {
  const cfgResp = await fetch("/assets/srch.json");
  const cfg = await cfgResp.json();
	window.cfg = cfg
  return cfg
};

async function initSrch() {
  const opts = await getCfg();
	cfgSrchClient(opts);
}

function adaptRequest(requests) {
}

function performSrch(requests) {
			let pp = {
				...opts,
				...requests[0].params,
			}
			let req = "?" + new URLSearchParams(pp).toString()
			let res = srch.search(req);
			let responses = JSON.parse(res)

			//responses.facetFields.forEach((facet) => {
				//let f = facet.items.map((item) => {
					//let i = {}
					//i[`${item.value}`] = item.count
				//console.log({item.value: item.count})
					//let i = {
						//item.label: item.count,
					//}
					//console.log(i)
					//return i
				//});
			//});

			//console.log(responses)
			return Promise.resolve({ results: [responses] });

};

const customSearchClient = {
	search: function (requests) {
		let pp = {
			...opts,
			...requests[0].params,
		}
		let req = "?" + new URLSearchParams(pp).toString()
		let res = srch.search(req);
		let responses = JSON.parse(res)

		//responses.facetFields.forEach((facet) => {
			//let f = facet.items.map((item) => {
				//let i = {}
				//i[`${item.value}`] = item.count
			//console.log({item.value: item.count})
				//let i = {
					//item.label: item.count,
				//}
				//console.log(i)
				//return i
			//});
		//});

		//console.log(responses)
		return Promise.resolve({ results: [responses] });
	}
};
