function sortings(cfg) {
	let sort = []
	cfg.sortableAttributes.forEach((by) => {
		let attr = by.split(":")[0];
		sort.push({
			value: `${attr}:desc`,
			label: `${attr} (desc)`,
		});
		sort.push({
			value: `${attr}:asc`,
			label: `${attr} (asc)`,
		});
	});
	return sort;
};

function facets(cfg) {
	return cfg.attributesForFaceting.map((attr) => {
		let op = 'or';
		const conj = (a) => a === attr;
		if (cfg.conjunctiveFacets.some(conj)) {
			op = "and"
		}
		return {
			attribute: attr,
			operator: op,
		};
	});
};


async function fetchWASM(url) {
	let response = await fetch(url);
	let wasm = await WebAssembly.instantiateStreaming(response, go.importObject);

	go.run(wasm.instance);
	return new Promise((r) => r(true));
};

// Get srch Options
async function fetchCfg(url) {
  const cfgResp = await fetch(url);
  const opts = await cfgResp.json();
  return opts
};

// Get data
async function fetchData(url) {
  const response = await fetch(url);
  const data = await response.json();
	return data
}

function cfgSrchClient(opts, data) {
  srch.newClient(queryStr(opts), JSON.stringify(data))
};

function queryStr(query) {
	let params = {
		...Alpine.store('srch').cfg,
		...query,
	}
	return "?" + new URLSearchParams(params).toString()
}

// adapt the instantsearch request
function adaptReq(requests) {
	if (requests[0].indexName !== "search") {
		let by = requests[0].indexName.split(":");
		requests[0].params.sortBy = by[0]
		requests[0].params.order = by[1]
	};

	let filters = requests[0].params.facetFilters
	if (filters) {
		requests[0].params.facetFilters = JSON.stringify(filters) 
	};

	return queryStr(requests[0].params)
}

// adapt the response to instantsearch format
function adaptRes(res) {
	let r = JSON.parse(res)
	//console.log(r.facets)
	let facetz = {};
	r.facetFields.forEach((facet) => {
		facetz[`${facet.attribute}`] = {}
			facet.items.forEach((item) => {
			facetz[`${facet.attribute}`][`${item.label}`] = item.count
		});
	});
	r.facets = facetz
	return r
}

let search = {};

// Start Search
async function initSearch(url) {
	//console.log("start instantsearch")
	const data = await fetchData(url);
	cfgSrchClient(Alpine.store('srch').cfg, data);

	// define custom client
	const customSearchClient = {
		search: function (requests) {
			//console.log(requests[0])
			let req = adaptReq(requests);
			let res = srch.search(req);
			let responses = adaptRes(res);
			return Promise.resolve({ results: [responses] });
		}
	};

	// set instantsearch options
	search = instantsearch({
		indexName: 'search',
		searchClient: customSearchClient,
		routing: {
			router: instantsearch.routers.history(),
		},
	});

	// add widgets
	search.addWidgets([
		sortBy({
			container: document.querySelector('#sort-by'),
			items: Alpine.store('srch').sortby,
			cssClasses: {
				select: ['form-select'],
				root: 'form-group',
			},
		}),
		customHits({
			container: document.querySelector('#hits'),
		}),
	]);

	Alpine.store('srch').facets.forEach((facet) => {
		search.addWidgets([
			customRefinementList({
				container: document.getElementById(`${facet.attribute}`),
				attribute: facet.attribute,
				//limit: 1000,
				operator: facet.operator,
				showMore: true,
				showMoreLimit: 20,
			})
		]);
	});

	// add more widgets
	search.addWidgets([
		customPagination({
			container: document.querySelector('#pagination'),
			totalPages: 5,
		}),
		clearRefinements({
			container: document.querySelector('#clear-refinements'),
			cssClasses: {
				button: ['btn', 'btn-success'],
			}, 
			templates: {
				resetLabel({ hasRefinements }, { html }) {
					return html`<span>reset</span>`;
				},
			},
		}),
		searchBox({
			container: document.querySelector('#searchbox'),
			cssClasses: {
				input: ['form-input', 'column', 'col-10'],
				form: 'columns m-1',
				submit: ['column', 'col-1', 'btn'],
				reset: ['column', 'col-1', 'btn'],
			}
		}),
	]);
	search.start();
}
