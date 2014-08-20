(function() {

    'use strict';

    //Create application
    var regoApp = angular.module('regoApp', ['ui.bootstrap']);

    //Create controller
    regoApp.controller('mainCtrl', ['$scope', '$http', '$location', '$modal', function($scope, $http, $location, $modal) {
        //Create emtpy matches
        $scope.matches = [];

        //Set default text
        $scope.regexpInput = 'r([a-z]+)go';
        $scope.stringInput = 'rego';
        $scope.findAllSubmatch = true;

        $scope.clearMatchResults = function() {
            //Clear all match output
            $scope.error = '';
            $scope.matches = [];
            $scope.matchResult = '';
        };

        $scope.clearAllFields = function() {
            $scope.regexpInput = '';
            $scope.stringInput = '';
            $scope.findAllSubmatch = true;
            $scope.clearMatchResults();
        };

        $scope.evaluateRegex = function() {
            //Retrieve updated regexp information
            var postData = {
                Expr: $scope.regexpInput,
                Text: $scope.stringInput,
                NumMatches: $scope.findAllSubmatch === 'true' || $scope.findAllSubmatch === true ? -1 : 1,
            };
            var uri = $scope.getBaseUrl() + "/eval_regexp/";
            $http.post(uri, postData)
                .success(function(data) {
                    //Clear results
                    $scope.clearMatchResults();

                    //Check for results
                    if (data.matches === null) {
                        return;
                    }

                    //Populate new results
                    var fullMatches = [];
                    for (var i = 0; i < data.matches.length; i++) {
                        var match = data.matches[i];

                        //Populate fullMatches list
                        fullMatches[i] = match[0];

                        for (var j = 1; j < match.length; j++) {
                            //Populate matches list
                            var groupName = data.groupsName.length >= j && data.groupsName[j - 1] !== '' ? data.groupsName[j - 1] : '-';
                            $scope.matches[j - 1] ={
                                    count: j - 1,
                                    groupName: groupName,
                                    matchText: match[j]
                            };
                        }
                    }
                    $scope.matchResult = fullMatches.join(" ");
                })
                .error(function(data) {
                    //Clear results
                    $scope.clearMatchResults();

                    //Populate error
                    $scope.error = data;
                });
        };

        $scope.shareRegex = function() {
            var postData = {
                Expr: $scope.regexpInput,
                Text: $scope.stringInput,
                NumMatches: $scope.findAllSubmatch === 'true' || $scope.findAllSubmatch === true ? -1 : 1,
            };
            var uri = $scope.getBaseUrl() + "/share_regexp/";
            $http.post(uri, postData)
            .success(function(data) {
                $modal.open({
                    templateUrl: 'shareModalContent.html',
                    controller: 'shareModalCtrl',
                    size: 'sm',
                    resolve: {
                        'shareUrl': function () {
                            return $scope.getBaseUrl() + "/load_regexp?key=" + data;
                        }
                    }
                });
            })
            .error(function(data) {
                $modal.open({
                    templateUrl: 'errorModalContent.html',
                    controller: 'errorModalCtrl',
                    resolve: {
                        'error': function () {
                            return data;
                        }
                    }
                });
            });
        };

        $scope.getBaseUrl = function() {
            return $location.protocol() + "://" + $location.host() + ":" + $location.port();
        };

        //Invoke evaluateRegex to display initial data to user
        $scope.evaluateRegex();
    }]);


    //Create share modal controller
    regoApp.controller('shareModalCtrl', ['$scope', 'shareUrl', function($scope, shareUrl) {
        $scope.shareUrl = shareUrl;
    }]);

    //Create error modal controller
    regoApp.controller('errorModalCtrl', ['$scope', 'error', function($scope, error) {
        $scope.error = error;
    }]);

})();
