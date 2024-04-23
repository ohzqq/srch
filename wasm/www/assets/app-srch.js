const { 
	connectHits,
	connectInfiniteHits,
	connectSearchBox,
	connectRefinementList,
	connectClearRefinements,
	connectSortBy,
	connectPagination,
} = instantsearch.connectors;

const { 
	hits,
	clearRefinements,
	pagination,
	searchBox,
	sortBy,
	refinementList,
} = instantsearch.widgets;

// Pagination
const renderPagination = (renderOptions, isFirstRender) => {
  const {
    pages,
    currentRefinement,
    nbPages,
    isFirstPage,
    isLastPage,
    refine,
    createURL,
  } = renderOptions;

  const container = document.querySelector('#pagination');

	container.innerHTML = `
    <div class="navbar text-center">
						<span class="navbar-section" style="padding:0 0.25em 0 0.25em;">
              <button
                href="${createURL(0)}"
								class="btn btn-primary text-center btn-block"
                data-value="${0}"
              >
                <<
              </button>
		</span>
						<span class="navbar-section" style="padding:0 0.25em 0 0.25em;">
              <button
                href="${createURL(currentRefinement - 1)}"
								class="btn btn-primary text-center btn-block"
                data-value="${currentRefinement - 1}"
              >
                <
              </button>
		</span>
      ${pages
        .map(
          page => `
						<span class="navbar-section" style="padding:0 0.25em 0 0.25em;">
              <button
                href="${createURL(page)}"
                data-value="${page}"
								class="btn btn-primary btn-block text-center ${currentRefinement === page ? 'active' : ''}"
              >
                ${page + 1}
              </button>
		</span>
          `
        )
        .join('')}
        ${`
						<span class="navbar-section" style="padding:0 0.25em 0 0.25em;">
						<button
							href="${createURL(currentRefinement + 1)}"
								class="btn btn-primary text-center btn-block"
							data-value="${currentRefinement + 1}"
						>
							>
						</button>
		</span>
						<span class="navbar-section" style="padding:0 0.25em 0 0.25em;">
						<button
							href="${createURL(nbPages - 1)}"
								class="btn btn-primary text-center btn-block"
							data-value="${nbPages - 1}"
						>
							>>
						</button>
		</span>
					`}
    </div>
  `;

	[...container.querySelectorAll('button')].forEach(element => {
    element.addEventListener('click', event => {
      refine(event.currentTarget.dataset.value);
    });
  });
};
const customPagination = connectPagination(renderPagination);

// SortBy
const renderSortBy = (renderOptions, isFirstRender) => {
  const {
    options,
    currentRefinement,
    refine,
    widgetParams,
    canRefine,
  } = renderOptions;

  if (isFirstRender) {
    const select = document.createElement('select');
	  select.setAttribute("label", "sort by")

    select.addEventListener('change', event => {
      refine(event.target.value);
    });

    widgetParams.container.appendChild(select);
  }

  const select = widgetParams.container.querySelector('select');

  select.disabled = !canRefine;

	let template = '';
	options.forEach((option) => {
		template += `
          <option
            value="${option.value}"
            ${option.value === currentRefinement ? 'selected' : ''}
          >
            ${option.label}
          </option>
		`;
	});
	select.innerHTML = template
};
const customSortBy = connectSortBy(renderSortBy);

// ClearRefinements
const renderClearRefinements = (renderOptions, isFirstRender) => {
	const { canRefine, refine, widgetParams } = renderOptions;
	if (isFirstRender) {
		const button = document.getElementById('clear-facets');
		button.addEventListener('click', () => {
			refine();
		});
	}
	widgetParams.container.querySelector('button').disabled = !canRefine;
};
const customClearRefinements = connectClearRefinements(renderClearRefinements);

// RefinementList
const renderRefinementList = (renderOptions, isFirstRender) => {
  const {
    items,
    isFromSearch,
    refine,
    createURL,
    isShowingMore,
    canToggleShowMore,
    searchForItems,
    toggleShowMore,
    widgetParams,
  } = renderOptions;

  if (isFirstRender) {
	  let button = document.getElementById(`show-more-${widgetParams.attribute}`)
    button.addEventListener('click', () => {
      toggleShowMore();
    });
  }

  widgetParams.container.querySelector('div').innerHTML = items
		//.filter((item) => item.count > cfg.minFacetCount)
    .map(
      item => `
        <label class="form-checkbox">
          <input type="checkbox"
						${item.isRefined ? 'checked' : ''}
            href="${createURL(item.value)}"
            data-value="${item.value}"
            style="font-weight: ${item.isRefined ? 'bold' : ''}"
          >
            <i class="form-icon"></i>${item.label} (${item.count})
          </a>
        </label>
      `
    )
    .join('');

  [...widgetParams.container.querySelectorAll('input')].forEach(element => {
    element.addEventListener('click', event => {
			//event.preventDefault();
      refine(event.currentTarget.dataset.value);
    });
  });

  let button = widgetParams.container.querySelector('button');

  button.disabled = !canToggleShowMore;
  button.textContent = isShowingMore ? 'Show less' : 'Show more';
};
const customRefinementList = connectRefinementList(renderRefinementList);

// SearchBox
const renderSearchBox = (renderOptions, isFirstRender) => {
  const { 
		query, 
		refine, 
		clear, 
		isSearchStalled, 
		widgetParams,
	} = renderOptions;

	let input = widgetParams.container.querySelector('#searchbox')

	input.addEventListener('input', event => {
		refine(event.target.value);
	});

	let button = widgetParams.container.querySelector('#clear-search')
	button.addEventListener('click', () => {
		clear();
	});

  widgetParams.container.querySelector('#searchbox').value = query;
};
const customSearchBox = connectSearchBox(renderSearchBox);


// Hits
const renderHits = (renderOptions, isFirstRender) => {
  const { 
		hits, 
		widgetParams,
	} = renderOptions;


	let template = '';
	hits.forEach((item, idx) => {
		let series = item.series ? `${item.series}, book ${item.series_index}` : ''
		template += `
			<div class="card">
				<div class="card-image">
					<!--<img src="${item.cover}"/>-->
				</div>
				<div class="card-header">
						<a href="${item.url}">
					<figure class="avatar avatar-xl float-left m-1" style="border-radius: 0;">
					<!--<img src="${item.cover}"/>-->
					</figure>
						</a>
					<div class="card-title">
						<a href="${item.url}">
						${item.title}
						</a>
					</div>
					<div class="card-subtitle">
						${series}
					</div>
				</div>
				<div class="card-footer">
				</div>
			</div>
		`;
	});
  widgetParams.container.innerHTML = template;
};
const customHits = connectHits(renderHits);

// Get srch Options
async function getCfg() {
  const cfgResp = await fetch("/assets/srch.json");
  const opts = await cfgResp.json();
	window.cfg = opts
  return opts
};

// Get data
async function getData() {
  const response = await fetch("/assets/data.json");
  const data = await response.json();
	return data
}

function cfgSrchClient(opts, data) {
	let params = new URLSearchParams(opts).toString()
  srch.newClient("?" + params, JSON.stringify(data))
};

// adapt the instantsearch request
function adaptReq(requests) {
	//console.log(requests[0].params);

	if (requests[0].indexName !== "search") {
		let by = requests[0].indexName.split(":");
		requests[0].params.sortBy = by[0]
		requests[0].params.order = by[1]
	};

	let filters = requests[0].params.facetFilters
	if (filters) {
		requests[0].params.facetFilters = JSON.stringify(filters) 
	};

	let pp = {
		...Alpine.store('srch').cfg,
		...requests[0].params,
	}
	return "?" + new URLSearchParams(pp).toString()
}

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
async function initSearch() {
	//console.log("start instantsearch")
	//const opts = await getCfg();
  const data = await getData();
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
