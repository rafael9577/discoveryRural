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
  backButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 1,
    borderColor: '#1B5E20',
    borderRadius: 12,
    backgroundColor: '#2E7D32',
  },

  suporteCartao: {
    width: "95%",
    flex: 1,
    alignSelf: "center",
    marginTop: 10,
  },
  containerPrincipal: {
    width: "100%",
    backgroundColor: "#2E7D32",
    borderRadius: 20,
    borderWidth: 2,
    borderColor: "#1B5E20",
    overflow: "hidden",
    paddingBottom: 15,
  },

  blocoSuperior: {
    flexDirection: "row",
    padding: 15,
  },
  imagemAnimal: {
    width: 100,
    height: 100,
    borderRadius: 12,
    borderWidth: 2,
    borderColor: "#1B5E20",
  },
  divisorVertical: {
    width: 2,
    backgroundColor: "#1B5E20",
    marginHorizontal: 12,
  },
  containerInformacoes: {
    flex: 1,
    justifyContent: "center",
  },
  nomeAnimal: {
    fontWeight: "bold",
    textTransform: "uppercase",
    fontSize: 16,
    color: "#FFFFFF",
  },
  divisorHorizontal: {
    height: 2,
    backgroundColor: "#1B5E20",
    marginVertical: 6,
  },
  nomeCientifico: {
    fontStyle: "italic",
    fontSize: 12,
    marginBottom: 6,
    color: "#E0E0E0",
  },
  linhaInformacao: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginTop: 4,
  },
  textoInformacao: {
    fontSize: 12,
    fontWeight: "600",
    color: "#E0E0E0",
  },
  valorInformacao: {
    fontSize: 12,
    fontWeight: "bold",
    color: "#FFFFFF",
  },

  blocoDescricao: {
    paddingHorizontal: 15,
    paddingTop: 15,
    borderTopWidth: 2,
    borderTopColor: "#1B5E20",
  },
  tituloDescricao2: {
    fontSize: 15,
    fontWeight: "bold",
    color: "#FFFFFF",
    marginBottom: 8,
  },
  linhaInformacao2: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginTop: 4,
  },
  textoInformacao2: {
    fontSize: 14,
    fontWeight: "600",
    color: "#E0E0E0",
  },
  valorInformacao2: {
    fontSize: 14,
    color: "#FFFFFF",
  },
  tituloDescricao: {
    fontSize: 15,
    fontWeight: "bold",
    color: "#E0E0E0",
    marginTop: 12,
    marginBottom: 4,
  },

  containerAviso: {
    backgroundColor: "rgba(255, 140, 0, 0.15)",
    borderRadius: 8,
    padding: 12,
    marginTop: 16,
    borderWidth: 1,
    borderColor: "rgba(255, 140, 0, 0.4)",
  },
  cabecalhoAviso: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 6,
  },
  iconeAviso: {
    marginRight: 6,
  },
  tituloAvisoDestaque: {
    fontWeight: "bold",
    fontSize: 16,
    color: "#FFB74D",
  },
  tituloAviso: {
    fontWeight: "bold",
    fontSize: 14,
    color: "#FFB74D",
    marginTop: 8,
  },
  textoAviso: {
    fontSize: 13,
    color: "#FFE0B2",
    marginTop: 2,
    marginBottom: 4,
  },

  modalContainer: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.9)",
    justifyContent: "center",
    alignItems: "center",
  },
  botaoFecharModal: {
    position: "absolute",
    top: 40,
    right: 20,
    zIndex: 10,
    padding: 10,
  },
  imagemModal: {
    width: "90%",
    height: "80%",
  },
});