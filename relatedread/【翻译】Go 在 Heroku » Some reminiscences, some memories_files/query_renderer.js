PopularQueryRenderer = function(container) {
  this.container = container;
};

PopularQueryRenderer.prototype.render = function(queries) {
  var numberOfPopularQueriesToShow = queries.popularQueries.length;
  if ((typeof queries.maxNumberOfPopularQueries == "number") &&
      (queries.maxNumberOfPopularQueries < numberOfPopularQueriesToShow)) {
    numberOfPopularQueriesToShow = queries.maxNumberOfPopularQueries;
  }
  for (var i = 0; i < numberOfPopularQueriesToShow; i++) {
    var a = document.createElement("a");
    a.setAttribute("href", queries.popularQueries[i].href);
    a.appendChild(document.createTextNode(queries.popularQueries[i].query));
    this.container.appendChild(a);
    this.container.appendChild(document.createTextNode(" "));
  }
};
