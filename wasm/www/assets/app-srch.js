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
		const button = document.createElement('ion-button');
		button.textContent = 'clear';

		button.addEventListener('click', () => {
			refine();
		});

		widgetParams.container.appendChild(button);
	}

	widgetParams.container.querySelector('ion-button').disabled = !canRefine;
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
		const ul = document.createElement('div');
		//ul.id = `list-${widgetParams.attribute}`
		ul.className = 'form-group'

    let button = document.createElement('button');
    button.textContent = 'Show more';
	  button.id = `show-more-${widgetParams.attribute}`
	  button.className = 'btn btn-primary'

    button.addEventListener('click', () => {
      toggleShowMore();
    });

		widgetParams.container.appendChild(ul);
    widgetParams.container.appendChild(button);
  }

  widgetParams.container.querySelector('div').innerHTML = items
		.filter((item) => item.count > cfg.minFacetCount)
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
	console.log(params)
};

function adaptReq(requests) {
	let pp = {
		...cfg,
		...requests[0].params,
	}
	return "?" + new URLSearchParams(pp).toString()
}

function adaptRes(res) {
	let r = JSON.parse(res)
	//console.log(r.facets)
	let facets = {};
	r.facetFields.forEach((facet) => {
		facets[`${facet.attribute}`] = {}
			facet.items.forEach((item) => {
			facets[`${facet.attribute}`][`${item.label}`] = item.count
		});

	});
	r.facets = facets
	console.log(r)

	return r
}

let search = {};

// Start Search
async function initSearch() {
	console.log("start instantsearch")
	
	const opts = await getCfg();
  const data = await getData();
	cfgSrchClient(opts, data);
	//

	console.log(opts)
	
	const customSearchClient = {
		search: function (requests) {
			console.log(requests[0])
			let req = adaptReq(requests);
			//console.log("request " + req)
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
	console.log("search")

	// add widgets
	search.addWidgets([
		//sortBy({
			//container: document.querySelector('#sort-by'),
			//items: Object.entries(opts.sortings).map((sort) => {
				//return {
					//value: sort[0],
					//label: `${sort[1].field} (${sort[1].order})`,
				//}
			//}),
			//cssClasses: {
				//select: ['form-select'],
				//root: 'form-group',
			//},
		//}),
		//hits({
		customHits({
			container: document.querySelector('#hits'),
			//transformItems(items) {
				//items.forEach((item) => {
					//let pd = {
						//title: item.title,
						//cover: item.cover,
						//feeds: [
							//{
								//type: "audio",
								//format: "aac",
								//url: item.url,
							//},
						//],
					//};
					//window[`podcastData${item.id}`] = JSON.stringify(pd);
				//});
				//return items;
			//},
		}),
	]);

	// Add refinementLists by aggregations
	//const facets = document.querySelector("#refinement-list");

	//console.log(opts.facets)

	//opts.attributesForFaceting.forEach((attr) => {
		//console.log(`'#${attr}'`);
		//let con = document.createElement("div");
		//con.id = attr;
		//facets.appendChild(con);

		//search.addWidgets([
			//customRefinementList({
				//container: con,
				//attribute: attr,
				//limit: 1000,
				//operator: facet.conjunction ? "and" : "or",
				//operator: "and",
				//showMore: true,
				//showMoreLimit: 20,
			//})
		//]);
	//});

	// add more widgets
	search.addWidgets([
		refinementList({
			container: document.querySelector('#tags'),
			attribute: "tags",
		}),
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

initSearch();
