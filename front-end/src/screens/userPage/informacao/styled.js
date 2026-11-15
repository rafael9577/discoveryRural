import { StyleSheet } from "react-native";

export default StyleSheet.create({
  safeArea: {
    backgroundColor: "#0A5133", 
    flex: 1,
  },
  
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: 'transparent',
    zIndex: 2, 
  },
  headerSpacer: {
    width: 40, 
  },
  headerTitleBox: {
    flex: 1,
    alignItems: 'center',
  },
  headerTitleText: {
    fontSize: 22,
    fontWeight: 'bold',
    color: '#FFFFFF', 
  },

  buscaContainer: {
    alignSelf: "center",
    width: "95%",
    marginTop: 10,
    marginBottom: 10,
  },
  buscaInputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#A7F3D0', 
    borderRadius: 20,
    borderWidth: 2,
    borderColor: '#1B5E20', 
    height: 55,
    paddingHorizontal: 15,
  },
  buscaIcon: {
    marginRight: 10,
  },
  buscaInput: {
    flex: 1,
    fontSize: 16,
    color: '#0A5133', 
  },

  cardContainer: {
    backgroundColor: "#2E7D32",
    flexDirection: "row",
    alignItems: "center",
    padding: 15,
    borderWidth: 2,
    borderColor: "#1B5E20",
    width: "95%",
    alignSelf: "center",
    marginBottom: 15, 
    borderRadius: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 3,
  },
  cardImage: {
    backgroundColor: "#1B5E20", 
    width: 90,
    height: 90,
    borderRadius: 15, 
  },
  cardTextArea: {
    flex: 1,
    justifyContent: "center",
    marginLeft: 15,
  },
  cardTitle: {
    fontWeight: "bold",
    fontSize: 18,
    color: "#FFFFFF",
    marginBottom: 5,
  },
  cardDescription: {
    fontSize: 13,
    color: "#E0E0E0", 
    lineHeight: 18,
  },

  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'transparent',
    marginTop: 50,
  },
  loadingText: {
    color: "#A7F3D0",
    marginTop: 10,
    fontSize: 16,
  },
  errorText: {
    color: "#FFFFFF",
    textAlign: "center",
    marginBottom: 15,
    fontSize: 16,
  },
  emptyText: {
    color: "#FFFFFF",
    textAlign: "center",
    fontSize: 16,
  },
  retryButton: {
    backgroundColor: "#268596",
    paddingVertical: 12,
    paddingHorizontal: 25,
    borderRadius: 20,
    borderWidth: 2,
    borderColor: '#1B5E20',
  },
  retryButtonText: {
    color: "#FFFFFF",
    fontWeight: "bold",
    fontSize: 16,
  },
});